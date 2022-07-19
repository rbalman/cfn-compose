package main

import (
	"fmt"
	"cfn-deploy/cfn"
	"cfn-deploy/workflow"
	"cfn-deploy/log"
	"os"
	"errors"
	"time"
	"context"
	"strconv"
)

var colors []string = []string{log.Blue, log.Yellow, log.Green, log.Magenta}
var logger = log.Logger{}
type Work struct {
	JobName string
	Job workflow.Job
	LogColor string
	DryRun bool
}

type Result struct {
	JobName string
	Error error
}

func main() {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	wf, err := workflow.Parse(os.Getenv("WORKFLOW"))
	if err != nil {
		logger.ColorPrintf(ctx, "[ERROR] Error while fetching workflow: %s\n", err.Error())
		os.Exit(1)
	}

	err = wf.Validate()
	if err  != nil {
		logger.ColorPrintf(ctx,"[ERROR] Failed while validating workflow: %s\n", err.Error())
		os.Exit(1)
	}

	work_ch := make(chan Work)
	results_ch := make(chan Result)

	//Generate the worker pool as pre the job counts
	jobCounts := len(wf.Jobs)
	for i := 0; i < jobCounts; i++ {
		go ExecuteJob(ctx, work_ch, results_ch, i)
	}

	queuedJob := 0
	var order uint = 0
	var dryRunFlag bool
	
	dryRunStr, ok := os.LookupEnv("DRY_RUN")
  if ok {
		dryRunFlag, err = strconv.ParseBool(dryRunStr)
		if err != nil {
			fmt.Printf("DRY_RUN should be either true/false %s", err)
			return
		}
	}

	//Publish work to the worker pool
	for {
		var jobsCountInOrder int = 0
		for name, job := range wf.Jobs {
			if job.Order == order {
				work_ch <- Work{JobName: name, Job: job, LogColor: colors[queuedJob], DryRun: dryRunFlag }
				queuedJob++
				jobsCountInOrder++
			}
		}
		fmt.Printf("[INFO] Dispatched Order: %d, JobCount: %d.\n", order, jobsCountInOrder)

		//wait for jobs for each order to complete
		for i := 0; i < jobsCountInOrder ; i++ {
			// logger.ColorPrintf(ctx,"[INFO] Order: %d, JobCount: %d. Waiting for result\n", order, jobsCountInOrder)
			r := <- results_ch
			// fmt.Printf("[DEBUG] Received Signal for Job: %s\n", r.JobName)
			if r.Error != nil {
				cancelCtx()
				fmt.Println("[INFO] Graceful wait for cancelled jobs")
				time.Sleep(time.Second * 10)
				// fmt.Printf("[DEBUG] Shutting down the process: %s\n", r.JobName)
				logger.Errorf("Workflow failed. Reason: %s", r.Error)
				return
			}
		}
			
		if queuedJob >= jobCounts {
			break
		}
		order++
	}

	cancelCtx()
	time.Sleep(time.Second*2)
	logger.ColorPrintf(ctx,"[INFO] Workflow Successfully executed!!")
}

func ExecuteJob(ctx context.Context, work_ch chan Work, results_ch chan Result, workerId int){
	defer func(){
		fmt.Printf("[DEBUG] Worker: %d exitting...\n", workerId)
	}()
	sess, err := getAWSSession(os.Getenv("AWS_PROFILE"), os.Getenv("AWS_REGION"))
	if err != nil {
		fmt.Printf("[ERROR] Error while getting Session: %s\n", err.Error())
		os.Exit(1)
	}
	
	cm := cfn.CFNManager{ Session: sess}

	for {
		select {
			case work := <- work_ch:
				//sleeping from readability
				time.Sleep(time.Millisecond * 500)
				name := work.JobName
				job := work.Job
				dryRun := work.DryRun
				jobCtx := context.WithValue(ctx, "logColor", work.LogColor)
				if dryRun {
					logger.ColorPrintf(jobCtx, "[INFO] DryRun started for Job: '%s'\n", name)
				}else{
					logger.ColorPrintf(jobCtx, "[INFO] Execution started for Job: '%s'\n", name)
				}

				for _, stack := range job.Stacks {
					var err error
					if dryRun {
						// logger.ColorPrintf(jobCtx,"[INFO] Executing DryRun on Job: '%s', Stack: '%s'\n", name, stack.StackName)
						err = stack.DryRun(jobCtx, cm)
					}else{
						logger.ColorPrintf(jobCtx,"[INFO] Applying Change for Job: '%s', Stack: '%s'\n", name, stack.StackName)
						err = stack.ApplyChanges(jobCtx, cm)
					}
					if err != nil {
						// logger.ColorPrintf(jobCtx, "[DEBUG] Job: %s Sending Fail Signal\n", name)
						errStr := fmt.Sprintf("[ERROR] Failed Job: '%s', Stack '%s', Error: %s\n", name, stack.StackName, err)
						logger.ColorPrintf(jobCtx, errStr)
						results_ch <- Result{
								Error: errors.New(errStr),
								JobName: name,
							}
						break;
					}
				}
				// logger.ColorPrintf(jobCtx, "[DEBUG] Sending Success Signal Job: %s\n", name)
				logger.ColorPrintf(jobCtx, "[INFO] Job: '%s' Completed Successfully!!\n", name)
				results_ch <- Result{JobName: name}

			case <- ctx.Done():
				// if err := ctx.Err(); err != nil {
				// 	fmt.Printf("[DEBUG] Cancel signal received Worker: %d, Info: %s\n", workerId, err)
				// }
				return
		}
	}
}
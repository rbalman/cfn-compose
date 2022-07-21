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
	// defer cancelCtx()

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

	var dryRunFlag bool = true
	dryRunStr, ok := os.LookupEnv("DRY_RUN")
  if ok {
		dryRunFlag, err = strconv.ParseBool(dryRunStr)
		if err != nil {
			fmt.Printf("DRY_RUN should be either true/false %s", err)
			return
		}
	}

	//Re-arrage jobs to ordered maps
	jobMap := make(map[uint][]workflow.Job)
	for name, job := range wf.Jobs {
		job.Name = name
		jobs := jobMap[job.Order]
		jobs = append(jobs, job)
		jobMap[job.Order] = jobs
	}

	workChan := make(chan Work)
	resultsChan := make(chan Result)
	//Generate the worker pool as pre the job counts
	jobCounts := len(wf.Jobs)
	for i := 0; i < jobCounts; i++ {
		go ExecuteJob(ctx, workChan, resultsChan, i)
	}


	//Dispatch Jobs in order
	for order, jobs := range jobMap {
		for index, job := range jobs {
			workChan <- Work{JobName: job.Name, Job: job, LogColor: colors[index], DryRun: dryRunFlag }
		}

		fmt.Printf("[INFO] Dispatched Order: %d, JobCount: %d.\n", order, len(jobs))

		//wait for jobs for each order to complete
		for i := 0; i < len(jobs) ; i++ {
			r := <- resultsChan
			if r.Error != nil {
				cancelCtx()
				fmt.Println("[INFO] Graceful wait for cancelled jobs")
				time.Sleep(time.Second * 10)
				logger.Errorf("Workflow failed. Reason: %s", r.Error)
				return
			}
		}
		fmt.Printf("[INFO] All Jobs completed for Run Order: %d\n\n", order)
	}

	cancelCtx()
	time.Sleep(time.Second*2)
	logger.ColorPrintf(ctx,"[INFO] Workflow Successfully Completed!!")
}


func ExecuteJob(ctx context.Context, workChan chan Work, resultsChan chan Result, workerId int){
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
			case work := <- workChan:
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
						resultsChan <- Result{
								Error: errors.New(errStr),
								JobName: name,
							}
						break;
					}
				}
				// logger.ColorPrintf(jobCtx, "[DEBUG] Sending Success Signal Job: %s\n", name)
				logger.ColorPrintf(jobCtx, "[INFO] Job: '%s' Completed Successfully!!\n", name)
				resultsChan <- Result{JobName: name}

			case <- ctx.Done():
				// if err := ctx.Err(); err != nil {
				// 	fmt.Printf("[DEBUG] Cancel signal received Worker: %d, Info: %s\n", workerId, err)
				// }
				return
		}
	}
}
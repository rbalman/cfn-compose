package main

import (
	"cfn-deploy/cfn"
	"cfn-deploy/logger"
	"cfn-deploy/workflow"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// var colors []string = []string{log.Blue, log.Yellow, log.Green, log.Magenta, log.Cyan}

type Work struct {
	JobName    string
	Job        workflow.Job
	DryRun     bool
	CfnManager cfn.CFNManager
}

type Result struct {
	JobName string
	Error   error
}

func main() {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	// defer cancelCtx()

	logger.Start(logger.DEBUG)

	wf, err := workflow.Parse(os.Getenv("WORKFLOW"))
	if err != nil {

		logger.Log.Errorf("Failed while fetching workflow: %s\n", err.Error())
		os.Exit(1)
	}

	err = wf.Validate()
	if err != nil {
		logger.Log.Errorf("Failed while validating workflow: %s\n", err.Error())
		os.Exit(1)
	}

	var dryRunFlag bool = true
	dryRunStr, ok := os.LookupEnv("DRY_RUN")
	if ok {
		dryRunFlag, err = strconv.ParseBool(dryRunStr)
		if err != nil {
			logger.Log.Errorf("DRY_RUN should be either true OR false. Error: %s", err)
			return
		}
	}

	//Re-arrange jobs to ordered maps
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

	// Exporting AWS_PROFILE and AWS_REGION for aws sdk client
	if val, ok := wf.Vars["AWS_PROFILE"]; ok {
		os.Setenv("AWS_PROFILE", val)
	}

	if val, ok := wf.Vars["AWS_REGION"]; ok {
		os.Setenv("AWS_REGION", val)
	}

	sess, err := getAWSSession()
	if err != nil {
		logger.Log.Errorf("Failed while creating AWS Session: %s\n", err.Error())
		os.Exit(1)
	}

	identity, err := getCallerIdentity(sess)
	if err != nil {
		logger.Log.Errorf("Failed to get AWS caller identity: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("##############################")
	fmt.Println("# Supplied AWS Configuration #")
	fmt.Println("##############################")
	printCallerIdentity(identity)
	fmt.Println()

	cm := cfn.CFNManager{Session: sess}

	//Dispatch Jobs in order
	for order, jobs := range jobMap {
		for _, job := range jobs {
			workChan <- Work{JobName: job.Name, Job: job, DryRun: dryRunFlag, CfnManager: cm}
		}

		logger.Log.Infof("Dispatched Order: %d, JobCount: %d.\n", order, len(jobs))

		//wait for jobs for each order to complete
		for i := 0; i < len(jobs); i++ {
			r := <-resultsChan
			if r.Error != nil {
				cancelCtx()
				logger.Log.Infoln("Graceful wait for cancelled jobs")
				time.Sleep(time.Second * 10)
				logger.Log.Errorf("Workflow failed. Error: %s", r.Error)
				return
			}
		}
		logger.Log.Infof("All Jobs completed for Dispatched Order: %d\n\n", order)
	}

	cancelCtx()
	time.Sleep(time.Second * 2)
	logger.Log.Infoln("Workflow Successfully Completed!!")
}

func ExecuteJob(ctx context.Context, workChan chan Work, resultsChan chan Result, workerId int) {
	defer func() {
		logger.Log.Debugf("[DEBUG] Worker: %d exiting...\n", workerId)
	}()

	for {
		select {
		case work := <-workChan:
			//sleeping from readability
			time.Sleep(time.Millisecond * 500)
			name := work.JobName
			job := work.Job
			dryRun := work.DryRun
			cm := work.CfnManager
			ctx := context.WithValue(ctx, "job", name)
			if dryRun {
				logger.Log.InfoCtxf(ctx, "DryRun started")
			} else {
				logger.Log.InfoCtxf(ctx, "Execution started")
			}

			for _, stack := range job.Stacks {
				ctx := context.WithValue(ctx, "stack", stack.StackName)
				var err error
				if dryRun {
					// logger.ColorPrintf(jobCtx,"[INFO] Executing DryRun on Job: '%s', Stack: '%s'\n", name, stack.StackName)
					err = stack.DryRun(ctx, cm)
				} else {
					// logger.Log.Infof("Applying Change for Job: '%s', Stack: '%s'\n", name, stack.StackName)
					logger.Log.InfoCtxf(ctx, "Applying Change")
					err = stack.ApplyChanges(ctx, cm)
				}

				if err != nil {
					// logger.ColorPrintf(jobCtx, "[DEBUG] Job: %s Sending Fail Signal\n", name)
					errStr := fmt.Sprintf("[JOB: %s] [STACK: %s]. Error: %s\n", name, stack.StackName, err)
					logger.Log.Infoln(errStr)
					// logger.ColorPrintf(jobCtx, errStr)
					resultsChan <- Result{
						Error:   errors.New(errStr),
						JobName: name,
					}
					break
				}
			}
			// logger.ColorPrintf(jobCtx, "[DEBUG] Sending Success Signal Job: %s\n", name)
			if dryRun {
				logger.Log.InfoCtxf(ctx, "DryRun Completed Successfully!!")
				// logger.Log.Infof("Job: '%s' DryRun Completed Successfully!!\n", name)
				// logger.ColorPrintf(jobCtx, "[INFO] Job: '%s' DryRun Completed Successfully!!\n", name)
			} else {
				logger.Log.InfoCtxf(ctx, "Completed Successfully!!")
				// logger.Log.Infof("Job: '%s' Completed Successfully!!\n", name)
				// logger.ColorPrintf(jobCtx, "[INFO] Job: '%s' Completed Successfully!!\n", name)
			}

			resultsChan <- Result{JobName: name}

		case <-ctx.Done():
			// if err := ctx.Err(); err != nil {
			// 	fmt.Printf("[DEBUG] Cancel signal received Worker: %d, Info: %s\n", workerId, err)
			// }
			return
		}
	}
}

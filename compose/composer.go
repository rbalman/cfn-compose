package compose

import (
	"github.com/balmanrawat/cfn-compose/cfn"
	"github.com/balmanrawat/cfn-compose/logger"
	"github.com/balmanrawat/cfn-compose/libs"
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

// var colors []string = []string{log.Blue, log.Yellow, log.Green, log.Magenta, log.Cyan}

type Work struct {
	JobName    string
	Job        Job
	DryRun     bool
	CfnManager cfn.CFNManager
}

type Result struct {
	JobName string
	Error   error
}

func Composer(configFile string, dryRun bool) {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	// defer cancelCtx()

	logger.Start(logger.INFO)

	wf, err := Parse(configFile)
	if err != nil {
		logger.Log.Errorf("Failed while fetching compose file: %s\n", err.Error())
		os.Exit(1)
	}

	err = wf.Validate()
	if err != nil {
		logger.Log.Errorf("Failed while validating compose file: %s\n", err.Error())
		os.Exit(1)
	}

	//MAP
	// key(order) value Job Array
	//
	//Re-arrange jobs to ordered maps
	jobMap := make(map[int][]Job)
	for name, job := range wf.Jobs {
		job.Name = name
		jobs, ok := jobMap[job.Order]
		if ok {
			jobs = append(jobs, job)
			jobMap[job.Order] = jobs
		} else {
			jobMap[job.Order] = []Job{job}
		}
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

	sess, err := libs.GetAWSSession()
	if err != nil {
		logger.Log.Errorf("Failed while creating AWS Session: %s\n", err.Error())
		os.Exit(1)
	}

	identity, err := libs.GetCallerIdentity(sess)
	if err != nil {
		logger.Log.Errorf("Failed to get AWS caller identity: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("##############################")
	fmt.Println("# Supplied AWS Configuration #")
	fmt.Println("##############################")
	libs.PrintCallerIdentity(identity)
	fmt.Println()

	cm := cfn.CFNManager{Session: sess}

	var order int
	var orders []int
	for key, _ := range jobMap {
		orders = append(orders, key)
	}

	logger.Log.Infof("TOTAL JOB COUNT: %d\n", jobCounts)
	sort.Ints(orders)
	//Dispatch Jobs in order
	for _, order = range orders {
		// for order, jobs := range jobMap {
		jobs, ok := jobMap[order]
		if !ok {
			continue
		}

		for _, job := range jobs {
			workChan <- Work{JobName: job.Name, Job: job, DryRun: dryRun, CfnManager: cm}
		}

		logger.Log.Infof("Dispatched Order: %d, JobCount: %d.\n", order, len(jobs))

		//wait for jobs for each order to complete
		for i := 0; i < len(jobs); i++ {
			r := <-resultsChan
			if r.Error != nil {
				cancelCtx()
				logger.Log.Infoln("Graceful wait for cancelled jobs")
				time.Sleep(time.Second * 10)
				logger.Log.Errorf("cfn compose failed. Error: %s", r.Error)
				return
			}
		}
		logger.Log.Infof("All Jobs completed for Dispatched Order: %d\n\n", order)
	}

	cancelCtx()
	time.Sleep(time.Second * 2)
	logger.Log.Infoln("Cfn Compose Successfully Completed!!")
}

func ExecuteJob(ctx context.Context, workChan chan Work, resultsChan chan Result, workerId int) {
	defer func() {
		logger.Log.Debugf("Worker: %d exiting...\n", workerId)
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
					err = stack.DryRun(ctx, cm)
				} else {
					logger.Log.InfoCtxf(ctx, "Applying Change")
					err = stack.ApplyChanges(ctx, cm)
				}

				if err != nil {
					errStr := fmt.Sprintf("[JOB: %s] [STACK: %s]. Error: %s\n", name, stack.StackName, err)
					logger.Log.Infoln(errStr)
					resultsChan <- Result{
						Error:   errors.New(errStr),
						JobName: name,
					}
					break
				}
			}

			if dryRun {
				logger.Log.InfoCtxf(ctx, "DryRun Completed Successfully!!")
			} else {
				logger.Log.InfoCtxf(ctx, "Completed Successfully!!")
			}

			resultsChan <- Result{JobName: name}

		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				logger.Log.DebugCtxf(ctx, "Cancel signal received Worker: %d, Info: %s\n", workerId, err)
			}
			return
		}
	}
}

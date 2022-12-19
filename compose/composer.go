package compose

import (
	"github.com/balmanrawat/cfn-compose/cfn"
	"github.com/balmanrawat/cfn-compose/logger"
	"github.com/balmanrawat/cfn-compose/libs"
	"github.com/balmanrawat/cfn-compose/config"
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

// var colors []string = []string{log.Blue, log.Yellow, log.Green, log.Magenta, log.Cyan}

type Work struct {
	Job        config.Job
	DryRun     bool
	DeployMode bool
	CfnManager cfn.CFNManager
}

type Result struct {
	JobName string
	Error   error
}

type Composer struct {
	Config config.ComposeConfig
	LogLevel string
	CherryPickedJob string
	DeployMode bool
	DryRun bool
	ConfigFile string
}

func (c *Composer)Apply() {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	cc, err := config.GetComposeConfig(c.ConfigFile)
	if err != nil {
		os.Exit(1)
	}

	err = cc.Validate()
	if err != nil {
		os.Exit(1)
	}

	c.Config = cc
	logger.StartWithLabel(c.LogLevel)

	var jobsMap map[int][]config.Job
	if c.CherryPickedJob != "" {
		jobsMap = cherryPickJob(c.CherryPickedJob, c.Config.Jobs)
		if len(jobsMap) == 0 {
			fmt.Printf("Err: Cannot find the selected job: %s in the config\n", c.CherryPickedJob)
		}
	}else{
		jobsMap = sortJobs(c.Config.Jobs)
	}
	
	orders := keys(jobsMap)
	if c.DeployMode {
		sort.Ints(orders)
	}else{
		sort.Sort(sort.Reverse(sort.IntSlice(orders)))
	}

	sess, err := libs.GetAWSSession()
	if err != nil {
		logger.Log.Errorf("Failed while creating AWS Session: %s\n", err.Error())
		os.Exit(1)
	}

	cm := cfn.CFNManager{Session: sess}

	workChan := make(chan Work)
	resultsChan := make(chan Result)
	//Generate the worker pool as pre the job counts
	for i := 0; i < len(c.Config.Jobs); i++ {
		go ExecuteJob(ctx, workChan, resultsChan, i)
	}
	logger.Log.Infof("TOTAL JOB COUNT: %d\n", len(c.Config.Jobs))
	
	//Dispatch Jobs in order
	for _, order := range orders {
		jobs, ok := jobsMap[order]
		if !ok {
			continue
		}

		for _, job := range jobs {
			workChan <- Work{Job: job, DryRun: c.DryRun, DeployMode: c.DeployMode, CfnManager: cm}
		}

		logger.Log.Infof("Dispatched Order: %d, JobCount: %d.\n", order, len(jobs))

		//wait for jobs in each order to complete
		for i := 0; i < len(jobs); i++ {
			r := <-resultsChan
			if r.Error != nil {
				cancelCtx()
				logger.Log.Infoln("Graceful wait for cancelled jobs")
				time.Sleep(time.Second * 10)
				logger.Log.Errorf("CFN compose failed. Error: %s", r.Error)
				return
			}
		}
		logger.Log.Infof("All Jobs completed for Dispatched Order: %d\n\n", order)
	}

	time.Sleep(time.Second * 2)
	logger.Log.Infoln("CFN Compose Successfully Completed!!")
}

func sortJobs(jobs map[string]config.Job) (map[int][]config.Job) {
	sortedJobs := make(map[int][]config.Job)
	for name, job := range jobs {
		job.Name = name

		jobs, ok := sortedJobs[job.Order]
		if ok {
			jobs = append(jobs, job)
			sortedJobs[job.Order] = jobs
		} else {
			sortedJobs[job.Order] = []config.Job{job}
		}
	}

	return sortedJobs
}

func keys(jobMap map[int][]config.Job) []int {
	var keys []int
	for key := range jobMap {
		keys = append(keys, key)
	}
	return keys
}

func cherryPickJob(jobName string, jobs map[string]config.Job) (map[int][]config.Job) {
	cherryPickedJob := make(map[int][]config.Job)
	for name, job := range jobs {
		if name == jobName {
			job.Name = name
			cherryPickedJob[job.Order] = []config.Job{job}
		}
	}
	return cherryPickedJob
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
			name := work.Job.Name
			job := work.Job
			dryRun := work.DryRun
			deployMode := work.DeployMode
			cm := work.CfnManager
			ctx := context.WithValue(ctx, "job", name)
			ctx = context.WithValue(ctx, "order", job.Order)

			if deployMode {
				for i:= 0 ;i < len(job.Stacks); i++{
					stack := job.Stacks[i]
					ctx := context.WithValue(ctx, "stack", stack.StackName)
					var err error
					if dryRun {
						err = stack.ApplyDryRun(ctx, cm)
					} else {
						logger.Log.InfoCtxf(ctx, "Applying Change...")
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
			}else {
				for i:= len(job.Stacks) - 1; i >= 0; i--{
					stack := job.Stacks[i]
					ctx := context.WithValue(ctx, "stack", stack.StackName)
					var err error
					if dryRun {
						err = stack.DestoryDryRun(ctx, cm)
					} else {
						logger.Log.InfoCtxf(ctx, "Destroying Stack...")
						err = stack.Destroy(ctx, cm)
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

func (c *Composer) Print() {
	fmt.Println("##########################")
	fmt.Println("# Compose Configuration #")
	fmt.Println("##########################")
	fmt.Printf("ConfigFile: %s\n", c.ConfigFile)
	if c.CherryPickedJob != "" {
		fmt.Printf("Selected Job: %s\n", c.CherryPickedJob)
	}
	fmt.Printf("DryRun: %t\n", c.DryRun)
	fmt.Printf("LogLevel: %s\n", c.LogLevel)
	fmt.Printf("DeployMode: %t\n\n", c.DeployMode)
}

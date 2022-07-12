package main

import (
	"fmt"
	"cfn-deploy/cfn"
	"cfn-deploy/workflow"
	"os"
	"errors"
)

func main() {
	wf, err := workflow.Parse(os.Getenv("WORKFLOW"))
	if err != nil {
		fmt.Printf("[ERROR] Error while fetching workflow: %s\n", err.Error())
		os.Exit(1)
	}

	err = wf.Validate()
	if err  != nil {
		fmt.Printf("[ERROR] Failed while validating workflow: %s\n", err.Error())
		os.Exit(1)
	}

	work_ch := make(chan map[string]workflow.Job)
	results_ch := make(chan error, 1)

	//Generate the worker pool as pre the job counts
	jobCounts := len(wf.Jobs)
	for i := 0; i < jobCounts ; i++ {
		go ExecuteJob(work_ch, results_ch, i)
	}

	queuedJob := 0
	var order uint = 0
	jobsCountInOrder := 0
	//Publish work to the worker pool
	for {
		for name, job := range wf.Jobs {
			if job.Order == order {
				work_ch <- map[string]workflow.Job{ name: job }
				queuedJob++
				jobsCountInOrder++
			}
		}

		//wait for jobs for each order to complete
		for i := 0; i < jobsCountInOrder ; i++ {
			fmt.Printf("[INFO] Waiting for result Order: %d, JobCount: %d\n", order, i)
			err := <- results_ch
			if err != nil {
				fmt.Sprintf("[ERROR] %s", err)
			}
		}
			
		if queuedJob >= jobCounts {
			break
		}
		order++
	}

	fmt.Println("[INFO] Workflow Successfully executed!!")
}


func ExecuteJob(work_ch chan map[string]workflow.Job, results_ch chan error, workerId int){
	sess, err := getAWSSession(os.Getenv("AWS_PROFILE"), os.Getenv("AWS_REGION"))
	if err != nil {
		fmt.Printf("[ERROR] Error while getting Session: %s\n", err.Error())
		os.Exit(1)
	}
	
	cm := cfn.CFNManager{ Session: sess}

	for jobMap := range work_ch {
		for name, job := range jobMap{
			for _, stack := range job.Stacks {
				fmt.Printf("[INFO] Applying Change for Order: %d, Job: %s, Stack: %s\n", job.Order, name, stack.StackName)
				err := stack.ApplyChanges(cm)
				fmt.Errorf("Failed while applying change for Job: %s, Stack %s, Error: %s\n", name, stack.StackName, err)
				results_ch <- errors.New(fmt.Sprintf("Failed while applying change for Job: %s, Stack %s, Error: %s\n", name, stack.StackName, err))
				fmt.Println()
			}
			fmt.Printf("[INFO] Job: Order: %d, %s Completed Successfully!!\n", job.Order, name)
		}
	}

	fmt.Printf("[INFO] Worker %d Retiring..\n", workerId)
}
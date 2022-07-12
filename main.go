package main

import (
	"fmt"
	"cfn-deploy/cfn"
	"cfn-deploy/workflow"
	"os"
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

	sess, err := getAWSSession(os.Getenv("AWS_PROFILE"), os.Getenv("AWS_REGION"))
	if err != nil {
		fmt.Printf("[ERROR] Error while getting Session: %s\n", err.Error())
		os.Exit(1)
	}
	
	cm := cfn.CFNManager{ Session: sess}

	for jobName, job := range wf.Jobs {
		for _, stack := range job.Stacks {
			fmt.Println("--------------------------------------------------")
			fmt.Printf("   Job: %s => Stack: %s\n", jobName, stack.StackName)
			fmt.Println("--------------------------------------------------")
			err = stack.ApplyChanges(cm)
			if err != nil {
				fmt.Printf("[ERROR] Failed while applying change for stack %s, Error: %s\n", stack.StackName, err)
				return
			}
			fmt.Println()
		}
	}
}
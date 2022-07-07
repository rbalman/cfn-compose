package main

import (
	"fmt"
	"cfn-deploy/cfn"
	"cfn-deploy/yml"
	"os"
)

func main() {
	workflow, err := yml.Parse(os.Getenv("WORKFLOW"))
	if err != nil {
		fmt.Printf("Error while fetching workflow: %s\n", err.Error())
	}

	sess, err := getAWSSession(os.Getenv("AWS_PROFILE"), os.Getenv("AWS_REGION"))
	if err != nil {
		fmt.Printf("Error while getting Session: %s\n", err.Error())
	}
	
	cm := cfn.CFNManager{ Session: sess}

	for jobName, job := range workflow.Jobs {
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
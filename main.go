package main

import (
	"fmt"
	"cfn-deploy/cfn"
	"cfn-deploy/yml"
)

func main() {
	workflow, err := yml.Parse("workflow.yml")
	if err != nil {
		fmt.Println("Error while fetching workflow: %s", err)
	}

	sess, err := getAWSSession("sputnik-pre-staging", "us-east-1")
	if err != nil {
		fmt.Println("Error while getting Session: %s", err)
	}
	
	cm := cfn.CFNManager{ Session: sess}

	for jobName, job := range workflow.Jobs {
		for _, stack := range job.Stacks {
			fmt.Printf("Job: %s => Stack: %s\n", jobName, stack.StackName)
			err = stack.ApplyChanges(cm)
			if err != nil {
				fmt.Printf("[ERROR] Failed while applying change for stack %s, Error: %s\n", stack.StackName, err)
				return
			}
			fmt.Println()
		}
	}
}
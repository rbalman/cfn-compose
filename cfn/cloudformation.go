package cfn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/session"
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"fmt"
	"time"
	"context"
	"cfn-deploy/log"
)

type CFNManager struct {
	Session *session.Session
}

var logger = log.Logger{}
//Details about status: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-describing-stacks.html

var CfnStatus []string = []string{"CREATE_IN_PROGRESS", "CREATE_COMPLETE", "ROLLBACK_IN_PROGRESS", "ROLLBACK_FAILED", "ROLLBACK_COMPLETE", "DELETE_IN_PROGRESS", "DELETE_FAILED", "UPDATE_IN_PROGRESS", "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS", "UPDATE_COMPLETE", "UPDATE_ROLLBACK_IN_PROGRESS", "UPDATE_ROLLBACK_FAILED", "UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS", "UPDATE_ROLLBACK_COMPLETE", "REVIEW_IN_PROGRESS"}

//////// MUTABLE OPERATIONS ////////
func (cm CFNManager) CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	svc := cloudformation.New(cm.Session)
	return svc.CreateStack(input)
}

func (cm CFNManager) UpdateStack(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	svc := cloudformation.New(cm.Session)
	return svc.UpdateStack(input)
}

func (cm CFNManager) DeleteStack(stackName string) (*cloudformation.DeleteStackOutput, error) {
	input := &cloudformation.DeleteStackInput{
		StackName:    aws.String(stackName),
	}

	svc := cloudformation.New(cm.Session)
	return svc.DeleteStack(input)
}

func (cm CFNManager) CreateChangeSet(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	svc := cloudformation.New(cm.Session)
	return svc.CreateChangeSet(input)
}

func (cm CFNManager) ExecuteChangeSet(ctx context.Context, input *cloudformation.ExecuteChangeSetInput) (*cloudformation.ExecuteChangeSetOutput, error) {
	svc := cloudformation.New(cm.Session)
	res, err := svc.ExecuteChangeSet(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				case "InvalidChangeSetStatus":
					logger.ColorPrintf(ctx,"[WARN] Change-set couldn't be applied. Warning: %s\n", err.Error())
				default:
					return res, err
			}
		}
	}

	return res, nil
	
}


//////// MUTABLE WAIT OPERATIONS ////////
func (cm CFNManager) CreateStackWithWait(ctx context.Context, input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	res, err := cm.CreateStack(input)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go loader(ch)
	err = cm.WaitStackCreateComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("stack create wait failed, ERROR: %s", err.Error()))
	}

	logger.ColorPrint(ctx,"\nWaiting Completed.")
	return res, nil
}


func (cm CFNManager) UpdateStackWithWait(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	res, err := cm.UpdateStack(input)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go loader(ch)
	err = cm.WaitStackUpdateComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("stack update wait failed, ERROR: %s", err.Error()))
	}

	return res, nil
}

func (cm CFNManager) DeleteStackWithWait(stackName string) (*cloudformation.DeleteStackOutput, error) {
	res, err := cm.DeleteStack(stackName)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go loader(ch)
	err = cm.WaitStackDeleteComplete(stackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("stack delete wait failed, ERROR: %s", err.Error()))
	}

	return res, nil
}

func (cm CFNManager) CreateChangeSetWithWait(input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	res, err := cm.CreateChangeSet(input)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go loader(ch)
	err = cm.WaitChangeSetCreateComplete(*input.StackName, *input.ChangeSetName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("stack change-set create failed, ERROR: %s", err.Error()))
	}

	return res, nil
}

func (cm CFNManager) ExecuteChangeSetWithWait(ctx context.Context, input *cloudformation.ExecuteChangeSetInput) (*cloudformation.ExecuteChangeSetOutput, error) {
	res, err := cm.ExecuteChangeSet(ctx, input)
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)
	ch := make(chan bool)
	go loader(ch)
	err = cm.WaitStackUpdateComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("stack update wait failed, ERROR: %s", err.Error()))
	}

	return res, nil
}


//////// WAIT OPERATIONS ////////
func (cm CFNManager) WaitStackCreateComplete(stackName string) error {
	svc := cloudformation.New(cm.Session)

	input := cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	// stack delete output is an empty struct
	err := svc.WaitUntilStackCreateComplete(&input)
	return err
}

func (cm CFNManager) WaitChangeSetCreateComplete(stackName string, changesetName string) error {
	svc := cloudformation.New(cm.Session)

	input := cloudformation.DescribeChangeSetInput{
		ChangeSetName: &changesetName,
		StackName: &stackName,
	}

	err := svc.WaitUntilChangeSetCreateComplete(&input)
	return err
}

func (cm CFNManager) WaitStackUpdateComplete(stackName string) error {
	svc := cloudformation.New(cm.Session)

	input := cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	err := svc.WaitUntilStackUpdateComplete(&input)
	return err
}

func (cm CFNManager) WaitStackDeleteComplete(stackName string) error {
	svc := cloudformation.New(cm.Session)

	input := cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	// stack delete output is an empty struct
	err := svc.WaitUntilStackDeleteComplete(&input)
	return err
}

//////// READ OPERATIONS ////////
func (cm CFNManager) DescribeStacks(stackName string) (*cloudformation.DescribeStacksOutput, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName:    aws.String(stackName),
	}

	svc := cloudformation.New(cm.Session)
	return svc.DescribeStacks(input)
}

func (cm CFNManager) ListStacks() (*cloudformation.ListStacksOutput, error) {
	svc := cloudformation.New(cm.Session)
	var pstatus []*string

	for _, s := range CfnStatus {
		//create copy of status
		cs := s
		pstatus = append(pstatus, &cs)
	}

	input := &cloudformation.ListStacksInput{
		StackStatusFilter: pstatus,
	}
	return svc.ListStacks(input)
}

func (cm CFNManager) DescribeChangeSet(stackName string, changeSetName string) (*cloudformation.DescribeChangeSetOutput, error) {
	svc := cloudformation.New(cm.Session)

	input := cloudformation.DescribeChangeSetInput{
		ChangeSetName: aws.String(changeSetName),
		StackName: &stackName,
	}

	return svc.DescribeChangeSet(&input)
}
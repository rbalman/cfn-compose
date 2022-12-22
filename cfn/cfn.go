package cfn

import (
	"github.com/rbalman/cfn-compose/logger"
	"github.com/rbalman/cfn-compose/libs"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CFNManager struct {
	Session *session.Session
}

//Details on CFN Status: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-describing-stacks.html
var CfnStatus []string = []string{"CREATE_COMPLETE", "UPDATE_COMPLETE", "ROLLBACK_COMPLETE",  "UPDATE_ROLLBACK_COMPLETE", "UPDATE_ROLLBACK_FAILED", "ROLLBACK_FAILED", "DELETE_FAILED", "CREATE_IN_PROGRESS","ROLLBACK_IN_PROGRESS", "DELETE_IN_PROGRESS", "UPDATE_IN_PROGRESS", "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS",  "UPDATE_ROLLBACK_IN_PROGRESS","UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS", "REVIEW_IN_PROGRESS"}

//////// MUTABLE OPERATIONS ////////
func (cm CFNManager) CreateStack(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	svc := cloudformation.New(cm.Session)
	return svc.CreateStack(input)
}

func (cm CFNManager) UpdateStack(input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	svc := cloudformation.New(cm.Session)
	return svc.UpdateStack(input)
}

func (cm CFNManager) DeleteStack(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
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
				logger.Log.WarnCtxf(ctx, "Change-set couldn't be applied. Warning: %s\n", err.Error())
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
	go libs.Loader(ctx, ch)
	err = cm.WaitStackCreateComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Wait CreateStack failed, ERROR: %s", err.Error()))
	}

	logger.Log.InfoCtxf(ctx, "Create Complete...")
	return res, nil
}

func (cm CFNManager) UpdateStackWithWait(ctx context.Context, input *cloudformation.UpdateStackInput) (*cloudformation.UpdateStackOutput, error) {
	res, err := cm.UpdateStack(input)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go libs.Loader(ctx, ch)
	err = cm.WaitStackUpdateComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Wait UpdateStack failed, ERROR: %s", err.Error()))
	}

	logger.Log.InfoCtxf(ctx, "Update Completed.")
	return res, nil
}

func (cm CFNManager) DeleteStackWithWait(ctx context.Context, input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
	res, err := cm.DeleteStack(input)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go libs.Loader(ctx, ch)
	err = cm.WaitStackDeleteComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Wait DeleteStack failed, ERROR: %s", err.Error()))
	}

	logger.Log.InfoCtxf(ctx, "Delete Complete...")
	return res, nil
}

func (cm CFNManager) CreateChangeSetWithWait(ctx context.Context, input *cloudformation.CreateChangeSetInput) (*cloudformation.CreateChangeSetOutput, error) {
	res, err := cm.CreateChangeSet(input)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool)
	go libs.Loader(ctx, ch)
	err = cm.WaitChangeSetCreateComplete(*input.StackName, *input.ChangeSetName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Change-set create failed, ERROR: %s", err.Error()))
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
	go libs.Loader(ctx, ch)
	err = cm.WaitStackUpdateComplete(*input.StackName)
	ch <- true
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Wait UpdateStack failed, ERROR: %s", err.Error()))
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
		StackName:     &stackName,
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
		StackName: aws.String(stackName),
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
		StackName:     &stackName,
	}

	return svc.DescribeChangeSet(&input)
}

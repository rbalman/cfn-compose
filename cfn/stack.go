package cfn

import (
	"context"
	"errors"
	"fmt"
	"github.com/rbalman/cfn-compose/libs"
	"github.com/rbalman/cfn-compose/logger"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const Regions string = "eu-north-1, ap-south-1, eu-west-3, eu-west-2, eu-west-1, ap-northeast-3,  ap-northeast-2, ap-northeast-1, sa-east-1, ca-central-1, ap-southeast-1, ap-southeast-2, eu-central-1, us-east-1, us-east-2, us-west-1, us-west-2"

type Stack struct {
	TemplateFile string            `yaml:"template_file,omitempty"`
	TemplateURL  string            `yaml:"template_url,omitempty"`
	StackName    string            `yaml:"stack_name"`
	Capabilities []string          `yaml:"capabilities,omitempty"`
	Parameters   map[string]string `yaml:"parameters,omitempty"`
	// ParametersFile   string            `yaml:"parameters_file"`
	Tags             map[string]string `yaml:"tags,omitempty"`
	TimeoutInMinutes int64             `yaml:"timeout,omitempty"`
	cm               CFNManager
}

/*
Stack is valid only when it satisfies all the below mentioned conditions:
- stack_name can't be empty
- one of template_url or template_file is mandatory, if both provided results into error
*/
func (s *Stack) Validate(index int) error {
	if s.StackName == "" {
		return fmt.Errorf("stack_name property for %d index stack is empty", index)
	}

	if s.TemplateFile == "" && s.TemplateURL == "" {
		return fmt.Errorf("one of the 'template_file' or 'template_url' property should be provided for %d index stack", index)
	}

	if s.TemplateFile != "" && s.TemplateURL != "" {
		return fmt.Errorf("can't provide value for both 'template_file' and 'template_url' property for %d index stack", index)
	}

	return nil
}

///TODO Make Create Input Methods DRY
func (s *Stack) createStackInput() (cloudformation.CreateStackInput, error) {
	var capabilities []*string
	for i := range s.Capabilities {
		capabilities = append(capabilities, &s.Capabilities[i])
	}

	var parameters []*cloudformation.Parameter
	for k, v := range s.Parameters {
		key := k
		value := v
		parameter := cloudformation.Parameter{
			ParameterKey:   &key,
			ParameterValue: &value,
		}
		parameters = append(parameters, &parameter)
	}

	var tags []*cloudformation.Tag
	for k, v := range s.Tags {
		key := k
		value := v
		tag := cloudformation.Tag{
			Key:   &key,
			Value: &value,
		}
		tags = append(tags, &tag)
	}

	input := cloudformation.CreateStackInput{
		Capabilities: capabilities,
		Parameters:   parameters,
		StackName:    &s.StackName,
		Tags:         tags,
	}

	if s.TemplateURL != "" {
		input.TemplateURL = &s.TemplateURL
	} else {
		templateBody, err := libs.ReadTemplate(s.TemplateFile)
		if err != nil {
			return cloudformation.CreateStackInput{}, err
		}
		input.TemplateBody = &templateBody
	}

	return input, nil
}

///TODO Make Update Input Methods DRY
func (s *Stack) updateStackInput() (cloudformation.UpdateStackInput, error) {
	var capabilities []*string
	for i := range s.Capabilities {
		capabilities = append(capabilities, &s.Capabilities[i])
	}

	var parameters []*cloudformation.Parameter
	for k, v := range s.Parameters {
		key := k
		value := v
		parameter := cloudformation.Parameter{
			ParameterKey:   &key,
			ParameterValue: &value,
		}
		parameters = append(parameters, &parameter)
	}

	var tags []*cloudformation.Tag
	for k, v := range s.Tags {
		key := k
		value := v
		tag := cloudformation.Tag{
			Key:   &key,
			Value: &value,
		}
		tags = append(tags, &tag)
	}

	input := cloudformation.UpdateStackInput{
		Capabilities: capabilities,
		Parameters:   parameters,
		StackName:    &s.StackName,
		Tags:         tags,
	}

	if s.TemplateURL != "" {
		input.TemplateURL = &s.TemplateURL
	} else {
		templateBody, err := libs.ReadTemplate(s.TemplateFile)
		if err != nil {
			return cloudformation.UpdateStackInput{}, err
		}
		input.TemplateBody = &templateBody
	}

	return input, nil
}

func (s *Stack) deleteStackInput() (cloudformation.DeleteStackInput, error) {
	return cloudformation.DeleteStackInput{
		StackName: &s.StackName,
	}, nil
}

///TODO Make Create Changeset Input Method DRY
func (s *Stack) createChangeSetInput(ctx context.Context) (cloudformation.CreateChangeSetInput, error) {
	var capabilities []*string
	for i := range s.Capabilities {
		capabilities = append(capabilities, &s.Capabilities[i])
	}

	var parameters []*cloudformation.Parameter
	for k, v := range s.Parameters {
		key := k
		value := v
		parameter := cloudformation.Parameter{
			ParameterKey:   &key,
			ParameterValue: &value,
		}
		parameters = append(parameters, &parameter)
	}

	var tags []*cloudformation.Tag
	for k, v := range s.Tags {
		key := k
		value := v
		tag := cloudformation.Tag{
			Key:   &key,
			Value: &value,
		}
		tags = append(tags, &tag)
	}

	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	changeSetName := s.StackName + "-" + nowStr
	includeNestedStacks := true

	input := cloudformation.CreateChangeSetInput{
		Capabilities:        capabilities,
		Parameters:          parameters,
		StackName:           &s.StackName,
		ChangeSetName:       &changeSetName,
		Tags:                tags,
		IncludeNestedStacks: &includeNestedStacks,
	}

	if s.TemplateURL != "" {
		input.TemplateURL = &s.TemplateURL
	} else {
		templateBody, err := libs.ReadTemplate(s.TemplateFile)
		if err != nil {
			return cloudformation.CreateChangeSetInput{}, err
		}
		input.TemplateBody = &templateBody
	}

	logger.Log.DebugCtxf(ctx, "Create Changeset Input %+v.\n", input)

	return input, nil
}

func (s *Stack) status(ctx context.Context, cm CFNManager) (string, error) {
	res, err := cm.DescribeStacks(s.StackName)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ValidationError":
				return "DOESN'T EXIST", nil
			default:
				// logger.ColorPrintf(ctx,"ERROR CODE: %s", aerr.Code())
				return "", errors.New(fmt.Sprintf("Failed while checking stack status, ERROR %+v", err.Error()))
			}
		}
	}

	cfnStack := res.Stacks[0]
	return *cfnStack.StackStatus, nil
}

func (s *Stack) ApplyChanges(ctx context.Context, cm CFNManager) error {
	status, err := s.status(ctx, cm)
	if err != nil {
		return err
	}

	switch status {
	case "DELETE_COMPLETE", "DOESN'T EXIST":
		logger.Log.InfoCtxf(ctx, "Creating Stack... as the stack is in %s state.\n", status)
		i, err := s.createStackInput()
		if err != nil {
			return err
		}

		_, err = cm.CreateStackWithWait(ctx, &i)
		if err != nil {
			return err
		}
	case "UPDATE_FAILED", "UPDATE_ROLLBACK_COMPLETE", "UPDATE_COMPLETE", "CREATE_COMPLETE":
		logger.Log.InfoCtxf(ctx, "Updating Stack... as the stack is in %s state.\n", status)
		i, err := s.updateStackInput()
		if err != nil {
			return err
		}

		_, err = cm.UpdateStackWithWait(ctx, &i)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "ValidationError":
					logger.Log.WarnCtxf(ctx, "Skipping... Update. Warning: %s\n", err.Error())
				default:
					return err
				}
			}
		}

	default:
		return errors.New(fmt.Sprintf("Stopping... the launch as the stack: %s status is: %s\n", s.StackName, status))
	}

	return nil
}

func (s *Stack) Destroy(ctx context.Context, cm CFNManager) error {
	status, err := s.status(ctx, cm)
	if err != nil {
		return err
	}

	switch status {
	case "DELETE_COMPLETE", "DOESN'T EXIST":
		logger.Log.InfoCtxf(ctx, "Skipping delete... as the stack is in %s state.\n", status)
		return nil

	case "CREATE_COMPLETE", "UPDATE_COMPLETE", "ROLLBACK_COMPLETE", "UPDATE_ROLLBACK_COMPLETE", "UPDATE_ROLLBACK_FAILED", "ROLLBACK_FAILED", "DELETE_FAILED":
		logger.Log.InfoCtxf(ctx, "Deleting Stack... as the stack is in %s state.\n", status)
		i, err := s.deleteStackInput()
		if err != nil {
			return err
		}
		_, err = cm.DeleteStackWithWait(ctx, &i)
		if err != nil {
			return err
		}

	default:
		return errors.New(fmt.Sprintf("Stopping... the deletion as the stack: %s status is: %s\n", s.StackName, status))
	}

	return nil
}

func (s *Stack) ApplyDryRun(ctx context.Context, cm CFNManager) error {
	status, err := s.status(ctx, cm)
	if err != nil {
		return err
	}

	switch status {
	case "DELETE_COMPLETE", "DOESN'T EXIST":
		logger.Log.InfoCtxf(ctx, "Stack Status: '%s'. Will be created.\n", status)

	case "UPDATE_FAILED", "UPDATE_ROLLBACK_COMPLETE", "UPDATE_COMPLETE", "CREATE_COMPLETE":
		// logger.ColorPrintf(ctx,"[DEBUG] Creating Changeset... for the stack: %s is at %s state\n", status, s.StackName)
		i, err := s.createChangeSetInput(ctx)
		if err != nil {
			return err
		}

		cs, err := cm.CreateChangeSetWithWait(ctx, &i)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotReady:") {
				logger.Log.InfoCtxf(ctx, "Stack Status: '%s'. No change detected.\n", status)
				return nil
			}
			return err
		}

		link := fmt.Sprintf("https://us-east-1.console.aws.amazon.com/cloudformation/home?region=%s#/stacks/changesets/changes?stackId=%s&changeSetId=%s", "us-east-1", url.QueryEscape(*cs.StackId), url.QueryEscape(*cs.Id))

		logger.Log.InfoCtxf(ctx, "Stack Status: '%s'. Will be updated.\n\tChangeSet Link: %s\n", status, link)

	default:
		logger.Log.InfoCtxf(ctx, "Can't run the operations as Stack is in %s state.\n", status)
	}

	return nil
}

func (s *Stack) DestroyDryRun(ctx context.Context, cm CFNManager) error {
	status, err := s.status(ctx, cm)
	if err != nil {
		return err
	}

	switch status {
	case "DELETE_COMPLETE", "DOESN'T EXIST":
		logger.Log.InfoCtxf(ctx, "Stack Status: '%s'. Delete will be Skipped.\n", status)

	case "CREATE_COMPLETE", "UPDATE_COMPLETE", "ROLLBACK_COMPLETE", "UPDATE_ROLLBACK_COMPLETE", "UPDATE_ROLLBACK_FAILED", "ROLLBACK_FAILED", "DELETE_FAILED":
		// logger.ColorPrintf(ctx,"[DEBUG] Creating Changeset... for the stack: %s is at %s state\n", status, s.StackName)
		logger.Log.InfoCtxf(ctx, "Stack Status: '%s'. Stack will be deleted.\n", status)
	default:
		logger.Log.InfoCtxf(ctx, "Can't run the operations as Stack is in %s state.\n", status)
	}

	return nil
}

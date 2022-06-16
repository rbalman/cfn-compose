package yml

import (
	"os"
	"gopkg.in/yaml.v2"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"cfn-deploy/cfn"
	"errors"
	"strconv"
	"strings"
	"time"
	"net/url"
)

const Regions string = "eu-north-1, ap-south-1, eu-west-3, eu-west-2, eu-west-1, ap-northeast-3,  ap-northeast-2, ap-northeast-1, sa-east-1, ca-central-1, ap-southeast-1, ap-southeast-2, eu-central-1, us-east-1, us-east-2, us-west-1, us-west-2"

type Stack struct {
	DependsOn string `yaml:"depends_on"`
	TemplateFile string `yaml:"template_file"`
	TemplateURL string `yaml:"template_url"`
	StackName string `yaml:"stack_name"`
	Capabilities []string `yaml:"capabilities"`
	Parameters map[string]string `yaml:"parameters"`
	ParametersFile string `yaml:"parameter_file"`
	Tags map[string]string `yaml:"tags"`
	TimeoutInMinutes int64 `yaml:"timeout"`
	EnableChangeSet bool `yaml:"enable_change_set"`
}

type Job struct {
	Description string `yml:"description"`
	Stacks map[string]Stack `yaml:"stacks`
	DependsOn string `yaml:"depends_on"`
}

type Workflow struct {
	Description string `yml:"description"`
	Jobs map[string]Job `yml:"jobs"`
}


//Parsing the configuration file
func Parse(file string) (Workflow, error) {
	var w Workflow
	data, err := os.ReadFile(file)
	if err != nil {
		return w, err
	}

	err = yaml.Unmarshal([]byte(data), &w)
	return w, err
}

///TODO Make Create Input Methods DRY
func (s *Stack)createStackInput() (cloudformation.CreateStackInput, error) {
	var capabilities []*string
	for _, capability := range s.Capabilities {
    capabilities = append(capabilities, &capability)
	}

	var parameters []*cloudformation.Parameter
	for k, v := range s.Parameters {
		key := k 
		value := v
		parameter := cloudformation.Parameter{
			ParameterKey: &key,
			ParameterValue: &value,
		}
		parameters = append(parameters, &parameter)
	}

	var tags []*cloudformation.Tag
	for k, v := range s.Tags {
		key := k
		value := v
		tag := cloudformation.Tag{
			Key: &key,
			Value: &value,
		}
		tags = append(tags, &tag)
	}

	templateBody, err := ReadTemplate(s.TemplateFile)
	if err != nil {
		return cloudformation.CreateStackInput{}, err
	}

	return cloudformation.CreateStackInput{
		Capabilities: capabilities,
		Parameters:   parameters,
		StackName:    &s.StackName,
		TemplateBody: &templateBody,
		Tags: tags,
	}, nil
}

///TODO Make Update Input Methods DRY
func (s *Stack)updateStackInput() (cloudformation.UpdateStackInput, error) {
	var capabilities []*string
	for _, capability := range s.Capabilities {
    capabilities = append(capabilities, &capability)
	}

	var parameters []*cloudformation.Parameter
	for k, v := range s.Parameters {
		key := k
		value := v
		parameter := cloudformation.Parameter{
			ParameterKey: &key,
			ParameterValue: &value,
		}
		parameters = append(parameters, &parameter)
	}
	
	var tags []*cloudformation.Tag
	for k, v := range s.Tags {
		key := k
		value := v
		tag := cloudformation.Tag{
			Key: &key,
			Value: &value,
		}
		tags = append(tags, &tag)
	}

	templateBody, err := ReadTemplate(s.TemplateFile)
	if err != nil {
		return cloudformation.UpdateStackInput{}, err
	}

	return cloudformation.UpdateStackInput{
		Capabilities: capabilities,
		Parameters:   parameters,
		StackName:    &s.StackName,
		TemplateBody: &templateBody,
		Tags: tags,
	}, nil
}


///TODO Make Create Changeset Input Method DRY
func (s *Stack)createChangeSetInput() (cloudformation.CreateChangeSetInput, error) {
	var capabilities []*string
	for _, capability := range s.Capabilities {
    capabilities = append(capabilities, &capability)
	}

	var parameters []*cloudformation.Parameter
	for k, v := range s.Parameters {
		key := k
		value := v
		parameter := cloudformation.Parameter{
			ParameterKey: &key,
			ParameterValue: &value,
		}
		parameters = append(parameters, &parameter)
	}
	
	var tags []*cloudformation.Tag
	for k, v := range s.Tags {
		key := k
		value := v
		tag := cloudformation.Tag{
			Key: &key,
			Value: &value,
		}
		tags = append(tags, &tag)
	}

	templateBody, err := ReadTemplate(s.TemplateFile)
	if err != nil {
		return cloudformation.CreateChangeSetInput{}, err
	}

	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	changeSetName := s.StackName + "-" + nowStr

	includeNestedStacks := true
	return cloudformation.CreateChangeSetInput{
		Capabilities: capabilities,
		Parameters:   parameters,
		StackName:    &s.StackName,
		TemplateBody: &templateBody,
		ChangeSetName: &changeSetName,
		Tags: tags,
		IncludeNestedStacks: &includeNestedStacks,
	}, nil
}


func (s *Stack) status(cm cfn.CFNManager) (string, error) {
	res, err := cm.DescribeStacks(s.StackName)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				case "ValidationError":
					return "DOESN'T EXIST", nil
				default:
					fmt.Println("ERROR CODE: ", aerr.Code())
					return "", errors.New(fmt.Sprintf("Failed while checking stack status, ERROR: %s", err.Error))
			}
		}
	}

	cfnStack := res.Stacks[0]
	return *cfnStack.StackStatus, nil
}

func (s *Stack) ApplyChanges(cm cfn.CFNManager) error {
	status, err := s.status(cm)
	if err != nil {
		return err
	}

	switch status {
		case "DELETE_COMPLETE", "DOESN'T EXIST":
			fmt.Printf("[INFO] Creating Stack... as the stack is in %s state\n", status)
			i, err := s.createStackInput()
			if err != nil {
				return err
			}

			_, err = cm.CreateStackWithWait(&i)
			if err != nil {
				return err
			}
		case "UPDATE_FAILED", "UPDATE_ROLLBACK_COMPLETE", "UPDATE_COMPLETE", "CREATE_COMPLETE":
			if s.EnableChangeSet {
				fmt.Printf("[INFO] Creating Changeset... as the enable_change_set flag is true, and the stack is at %s state\n", status)
				i, err := s.createChangeSetInput()
				if err != nil {
					return err
				}

				cs, err := cm.CreateChangeSetWithWait(&i)
				if err != nil {
					if strings.Contains(err.Error(), "ResourceNotReady:") {
						fmt.Println("[INFO] Skipping.. No changes found while creating change-set")
						return nil
					}
					return err
				}

				link := fmt.Sprintf("https://us-east-1.console.aws.amazon.com/cloudformation/home?region=%s#/stacks/changesets/changes?stackId=%s&changeSetId=%s", "us-east-1", url.QueryEscape(*cs.StackId), url.QueryEscape(*cs.Id))
				fmt.Printf("\n[INFO] Review Required\n Link: %s \n", link)
				fmt.Println()
				fmt.Printf("[INPUT] Execute ChangeSet(y|Y) or Discard ChangeSet(n|N)?:")
				var execute string
				fmt.Scanln(&execute)

				if execute == "y" || execute == "Y" {
					fmt.Printf("[INFO] Executing the ChangeSet as input was: %s\n", execute)
					_, err = cm.ExecuteChangeSetWithWait(&cloudformation.ExecuteChangeSetInput{
						ChangeSetName: cs.Id,
						StackName: cs.StackId,
					})
				}else {
					fmt.Printf("[INFO] Skipping the ChangeSet as input was: %s\n", execute)
				}
				return nil
			}


			fmt.Printf("[INFO] Updating Stack... as the stack is in %s state\n", status)
			i, err := s.updateStackInput()
			if err != nil {
				return err
			}

			_, err = cm.UpdateStackWithWait(&i)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
						case "ValidationError":
							fmt.Printf("[WARN] Skipping... Update. Warning: %s\n", err.Error())
						default:
							return err
					}
				}
			}

		default:
			return errors.New(fmt.Sprintf("Stopping... the launch as the stack status is in %s state\n", status))
	}

	return nil
}

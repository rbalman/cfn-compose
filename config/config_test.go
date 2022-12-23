package config

import (
	"fmt"
	"github.com/rbalman/cfn-compose/cfn"
	"testing"
	"strconv"
)

func TestValidateComposeConfig(t *testing.T) {
	t.Log("When There are no flows in Compose")
	{
		var cc ComposeConfig
		err := cc.Validate()
		if err == nil {
			t.Fatal(fmt.Sprintf("Validation should not return nil"))
		}
	}

	t.Log("When flows count is above the limit")
	{
		cc := ComposeConfig{
			Flows: generateFlowsMap(flowCountLimit +1, 1),
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When flow order is negative")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Order: -1,
					Stacks: []cfn.Stack{
						{},
					},
				},
				"flow2": {
					Order: 1,
					Stacks: []cfn.Stack{
						{},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When flow order is greater than 100")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Order: 101,
					Stacks: []cfn.Stack{
						{},
					},
				},
				"flow2": {
					Order: 1,
					Stacks: []cfn.Stack{
						{},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When one or more flow doesn't have any stack")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Stacks: []cfn.Stack{
						{},
						{},
					},
				},
				"flow2": {},
				"flow6": {},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When one or more flow has stacks above the limit")
	{
		cc := ComposeConfig{
			Flows: generateFlowsMap(1, stackCountLimit + 1),
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When stack doesn't have a name")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Stacks: []cfn.Stack{
						{},
					},
				},
				"flow2": {
					Stacks: []cfn.Stack{
						{},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When stack doesn't provide both template_url/template_file name")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Stacks: []cfn.Stack{
						{
							StackName: "s1-stack",
						},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}

	}

	t.Log("When stack both template_url/template_file is provided")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Stacks: []cfn.Stack{
						{
							StackName:    "s1-stack",
							TemplateFile: "/Users/mockuser/cfn-templates/template.yaml",
							TemplateURL:  "https://artifactory.amazonaws.com/cfn-templates/template.yaml",
						},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When only template_url is provided")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Stacks: []cfn.Stack{
						{
							StackName:   "s1-stack",
							TemplateURL: "https://artifactory.amazonaws.com/cfn-templates/template.yaml",
						},
					},
				},
			},
		}

		err := cc.Validate()
		if err != nil {
			t.Fatal(fmt.Sprintf("Validation should return nil but found error: %s", err))
		}
	}

	t.Log("When only template_file is provided")
	{
		cc := ComposeConfig{
			Flows: map[string]Flow{
				"flow1": {
					Stacks: []cfn.Stack{
						{
							StackName:    "s1-stack",
							TemplateFile: "/Users/mockuser/cfn-templates/template.yaml",
						},
					},
				},
			},
		}

		err := cc.Validate()
		if err != nil {
			t.Fatal(fmt.Sprintf("Validation should return nil but found error: %s", err))
		}
	}
}

func generateFlowsMap(flowCount, stackCount int) map[string]Flow {
	m := make(map[string]Flow)
	var stacks []cfn.Stack
	
	for i := 0; i < stackCount ; i++ {
		name := "stack" + strconv.Itoa(i)
		stacks = append(stacks, cfn.Stack{StackName:name, TemplateFile: name + ".yml"  })
	}

	for i := 0; i < flowCount ; i++ {
		key := "flow" + strconv.Itoa(i)
		m[key] = Flow{
			Description: key, 
			Stacks: stacks,
		}
	}

	return m
}
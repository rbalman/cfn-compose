package config

import (
	"fmt"
	"testing"
	"github.com/rbalman/cfn-compose/cfn"
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
			Flows: map[string]Flow{
				"flow1": Flow{},
				"flow2": Flow{},
				"flow3": Flow{},
				"flow4": Flow{},
				"flow5": Flow{},
				"flow6": Flow{},
			},
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
				"flow1": Flow{
					Order: -1,
					Stacks: []cfn.Stack{
						cfn.Stack{},
					},
				},
				"flow2": Flow{
					Order: 1,
					Stacks: []cfn.Stack{
						cfn.Stack{},
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
				"flow1": Flow{
					Order: 101,
					Stacks: []cfn.Stack{
						cfn.Stack{},
					},
				},
				"flow2": Flow{
					Order: 1,
					Stacks: []cfn.Stack{
						cfn.Stack{},
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
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{},
						cfn.Stack{},
					},
				},
				"flow2": Flow{},
				"flow6": Flow{},
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
			Flows: map[string]Flow{
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
						cfn.Stack{},
					},
				},
				"flow2": Flow{},
				"flow3": Flow{},
			},
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
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{},
					},
				},
				"flow2": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{},
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
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{
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
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{
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
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{
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
				"flow1": Flow{
					Stacks: []cfn.Stack{
						cfn.Stack{
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

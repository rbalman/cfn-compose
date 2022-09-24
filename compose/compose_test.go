package compose

import (
	"fmt"
	"testing"
)

func TestValidateComposeConfig(t *testing.T) {
	t.Log("When There are no jobs in Compose")
	{
		var cc ComposeConfig
		err := cc.Validate()
		if err == nil {
			t.Fatal(fmt.Sprintf("Validation should not return nil"))
		}
	}

	t.Log("When jobs count is above the limit")
	{
		cc := ComposeConfig{
			Jobs: map[string]Job{
				"job1": Job{},
				"job2": Job{},
				"job3": Job{},
				"job4": Job{},
				"job5": Job{},
				"job6": Job{},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When job order is negative")
	{
		cc := ComposeConfig{
			Jobs: map[string]Job{
				"job1": Job{
					Order: -1,
					Stacks: []Stack{
						Stack{},
					},
				},
				"job2": Job{
					Order: 1,
					Stacks: []Stack{
						Stack{},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When job order is greater than 100")
	{
		cc := ComposeConfig{
			Jobs: map[string]Job{
				"job1": Job{
					Order: 101,
					Stacks: []Stack{
						Stack{},
					},
				},
				"job2": Job{
					Order: 1,
					Stacks: []Stack{
						Stack{},
					},
				},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When one or more job doesn't have any stack")
	{
		cc := ComposeConfig{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{},
						Stack{},
					},
				},
				"job2": Job{},
				"job6": Job{},
			},
		}

		err := cc.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When one or more job has stacks above the limit")
	{
		cc := ComposeConfig{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
						Stack{},
					},
				},
				"job2": Job{},
				"job3": Job{},
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
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{},
					},
				},
				"job2": Job{
					Stacks: []Stack{
						Stack{},
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
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{
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
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{
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
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{
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
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: []Stack{
						Stack{
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

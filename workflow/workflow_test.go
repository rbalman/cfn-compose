package workflow

import (
	"fmt"
	"testing"
)

func TestValidateWorkflow(t *testing.T) {
	t.Log("When There are no jobs in Workflow")
	{
		var w Workflow
		err := w.Validate()
		if err != nil {
			t.Fatal(fmt.Sprintf("Validation should return nil found %s", err))
		}
	}

	t.Log("When jobs count is above the limit")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{},
				"job2": Job{},
				"job3": Job{},
				"job4": Job{},
				"job5": Job{},
				"job6": Job{},
			},
		}

		err := w.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When one or more job doesn't have any stack")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1": Stack{},
						"s2": Stack{},
					},
				},
				"job2": Job{},
				"job6": Job{},
			},
		}

		err := w.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When one or more job has stacks above the limit")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1":  Stack{},
						"s2":  Stack{},
						"s3":  Stack{},
						"s4":  Stack{},
						"s5":  Stack{},
						"s6":  Stack{},
						"s7":  Stack{},
						"s8":  Stack{},
						"s9":  Stack{},
						"s10": Stack{},
						"s11": Stack{},
						"s12": Stack{},
						"s13": Stack{},
						"s14": Stack{},
						"s15": Stack{},
						"s16": Stack{},
						"s17": Stack{},
						"s18": Stack{},
						"s19": Stack{},
						"s20": Stack{},
						"s21": Stack{},
						"s22": Stack{},
						"s23": Stack{},
						"s24": Stack{},
						"s25": Stack{},
						"s26": Stack{},
						"s27": Stack{},
						"s28": Stack{},
						"s29": Stack{},
						"s30": Stack{},
						"s31": Stack{},
						"s32": Stack{},
						"s33": Stack{},
						"s34": Stack{},
						"s35": Stack{},
						"s36": Stack{},
					},
				},
				"job2": Job{},
				"job3": Job{},
			},
		}

		err := w.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When stack doesn't have a name")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1": Stack{},
					},
				},
				"job2": Job{
					Stacks: map[string]Stack{
						"s1": Stack{},
					},
				},
			},
		}

		err := w.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When stack has doesn't provide both template_url/template_file name")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1": Stack{
							StackName: "s1-stack",
						},
					},
				},
			},
		}

		err := w.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When stack both template_url/template_file is provided")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1": Stack{
							StackName:    "s1-stack",
							TemplateFile: "/Users/mockuser/cfn-templates/template.yaml",
							TemplateURL:  "https://artifactory.amazonaws.com/cfn-templates/template.yaml",
						},
					},
				},
			},
		}

		err := w.Validate()
		if err == nil {
			t.Fatal("Validation should return error but found nil", err)
		}
	}

	t.Log("When only template_url is provided")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1": Stack{
							StackName:   "s1-stack",
							TemplateURL: "https://artifactory.amazonaws.com/cfn-templates/template.yaml",
						},
					},
				},
			},
		}

		err := w.Validate()
		if err != nil {
			t.Fatal(fmt.Sprintf("Validation should return nil but found error: %s", err))
		}
	}

	t.Log("When only template_file is provided")
	{
		w := Workflow{
			Jobs: map[string]Job{
				"job1": Job{
					Stacks: map[string]Stack{
						"s1": Stack{
							StackName:    "s1-stack",
							TemplateFile: "/Users/mockuser/cfn-templates/template.yaml",
						},
					},
				},
			},
		}

		err := w.Validate()
		if err != nil {
			t.Fatal(fmt.Sprintf("Validation should return nil but found error: %s", err))
		}
	}
}

func TestPrepareVariables(t *testing.T) {
	t.Log("When There are no jobs in Workflow")
	{
		var w Workflow
		err := w.Validate()
		if err != nil {
			t.Fatal(fmt.Sprintf("Validation should return nil found %s", err))
		}
	}
}

package workflow

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// type Workflow struct {
// 	Description string            `yml:"description" json:"description"`
// 	Jobs        map[string]Job    `yml:"jobs" json:"jobs"`
// 	Vars        map[string]string `yml:"vars" json:"vars"`
// }

type Workflow struct {
	Description string            `yaml:"description"`
	Jobs        map[string]Job    `yaml:"jobs"`
	Vars        map[string]string `yaml:"vars"`
}

var jobCountLimit int = 5

// type Job struct {
// 	Name        string  `yml:"name" json:"name"`
// 	Description string  `yml:"description" json:"description"`
// 	Stacks      []Stack `yaml:"stacks json:"stacks"`
// 	Order       uint    `yaml:"order" json:"order"`
// }

type Job struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Stacks      []Stack `yaml:"stacks"`
	Order       uint    `yaml:"order"`
}

/*
Job is valid when all of the below conditions are true:
- When Stack counts is <= stackCountLimit
- order property should be a valid unsigned integer
- When all stacks are valid
*/
func (j *Job) Validate(name string) error {
	if len(j.Stacks) >= stackCountLimit || len(j.Stacks) == 0 {
		return fmt.Errorf("Stack count is %d for Job: %s, should be '> 0 and <= %d'", len(j.Stacks), name, stackCountLimit)
	}

	for i, stack := range j.Stacks {
		err := stack.Validate(i)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Workflow is valid when all of the below conditions are true:
- When Job counts is <= jobCoutLimit
- When all jobs are valid
- When all stacks inside the jobs are valid
*/
func (w *Workflow) Validate() error {
	if len(w.Jobs) > jobCountLimit {
		return fmt.Errorf("Job count is %d, should be <= %d", len(w.Jobs), jobCountLimit)
	}

	for jname, job := range w.Jobs {
		if err := job.Validate(jname); err != nil {
			return err
		}
	}

	return nil
}

//Parsing the configuration file
func Parse(file string) (Workflow, error) {
	var w Workflow
	// var wf Workflow
	data, err := os.ReadFile(file)
	if err != nil {
		return w, err
	}

	vars, err := prepareVariables(data)
	if err != nil {
		return w, err
	}

	t, err := template.New("WorkFlowTemplate").Funcs(template.FuncMap{
		"shell": func(bin string, args ...string) string {
			if len(bin) < 1 {
				log.Fatal(errors.New("shell function requires at least one argument, which must be name of the binary."))
				return ""
			}
			cmd := exec.Command(bin, args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				log.Fatal(errors.New(fmt.Sprintf("Shell command failed with: %s\n Error: %s", stderr.String(), err)))
				return ""
			}

			// log.Println("OUTPUT: ", stdout.String())
			return strings.TrimSpace(stdout.String())
		},
	}).Parse(string(data))
	if err != nil {
		return w, err
	}

	final_template_file, err := os.Create(os.Getenv("WORKFLOW") + ".final.template")
	if err != nil {
		return w, err
	}
	defer final_template_file.Close()
	t.Execute(final_template_file, vars)

	varsData, err := os.ReadFile(os.Getenv("WORKFLOW") + ".final.template")
	if err != nil {
		return w, err
	}

	err = yaml.Unmarshal([]byte(varsData), &w)
	if err != nil {
		return w, err
	}

	w.Vars = vars

	return w, err
}

func prepareVariables(data []byte) (map[string]string, error) {
	varData := struct {
		Vars map[string]string `yml:"vars"`
	}{}

	varStruct := struct {
		Vars map[string]string `yml:"vars"`
	}{}

	err := yaml.Unmarshal(data, &varData)
	if err != nil {
		return varStruct.Vars, err
	}

	varBytes, err := yaml.Marshal(varData)
	if err != nil {
		return varStruct.Vars, err
	}

	t, err := template.New("VarsTemplate").Funcs(template.FuncMap{
		"shell": func(bin string, args ...string) string {
			if len(bin) < 1 {
				log.Fatal(errors.New("shell function requires at least one argument, which must be name of the binary."))
				return ""
			}
			cmd := exec.Command(bin, args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				log.Fatal(errors.New(fmt.Sprintf("Shell command failed with: %s\n Error: %s", stderr.String(), err)))
				return ""
			}

			// log.Println("OUTPUT: ", stdout.String())
			return strings.TrimSpace(stdout.String())
		},
	}).Parse(string(varBytes))

	if err != nil {
		return varStruct.Vars, err
	}

	vars_file, err := os.Create(os.Getenv("WORKFLOW") + ".vars.template")
	if err != nil {
		return varStruct.Vars, err
	}

	defer vars_file.Close()
	err = t.Execute(vars_file, nil)
	if err != nil {
		return varStruct.Vars, err
	}

	varFileData, err := os.ReadFile(os.Getenv("WORKFLOW") + ".vars.template")
	if err != nil {
		return varStruct.Vars, err
	}

	err = yaml.Unmarshal([]byte(varFileData), &varStruct)
	if err != nil {
		return varStruct.Vars, err
	}

	err = overrideWithEnvs(varStruct.Vars)
	return varStruct.Vars, err
}

func overrideWithEnvs(varsMap map[string]string) error {
	var err error
	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		varsMap[split_v[0]] = split_v[1]
	}

	return err
}

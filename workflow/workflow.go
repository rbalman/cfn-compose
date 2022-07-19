package workflow

import (
	"os"
	"gopkg.in/yaml.v2"
	"fmt"
	"errors"
	"strings"
	"text/template"
	"os/exec"
	"bytes"
	"log"
)


type Workflow struct {
	Description string `yml:"description"`
	Jobs map[string]Job `yml:"jobs"`
	Vars map[string]string `yml:"vars"`
}

var jobCountLimit int = 5
type Job struct {
	Description string `yml:"description"`
	Stacks map[string]Stack `yaml:"stacks`
	Order uint `yaml:"order"`
}

/*
Job is valid when all of the below conditions are true:
- When Stack counts is <= stackCountLimit
- order property should be a valid unsigned integer
- When all stacks are valid
*/
func (j *Job) Validate(name string) error {
	if len(j.Stacks) >= stackCountLimit || len(j.Stacks) == 0 {
		return fmt.Errorf("Stack count is %d for Job: %s, should be '> 0 and <= %d'", len(j.Stacks), name,  stackCountLimit)
	}

	for name, stack := range j.Stacks {
		err := stack.Validate(name);
		if err  != nil {
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
	if len(w.Jobs) >= jobCountLimit {
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
	var wf Workflow
	data, err := os.ReadFile(file)
	if err != nil {
		return w, err
	}
	err = yaml.Unmarshal([]byte(data), &wf)
	if err != nil {
		return w, err
	}

	err = overrideWithEnvs(wf.Vars)
	if err != nil {
		return w, err
	}

	t, err := template.New("WorkflowTemplate").Funcs(template.FuncMap{
    "shell": func(bin string, args ...string) string {
			if len(bin) < 1 {
				log.Fatal(errors.New("Shell command requires at least one argument, which must be name of the binary."))
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

	tfile, err := os.Create("/tmp/workflow.yaml")
	if err != nil {
		return w, err
	}
	defer tfile.Close()
	t.Execute(tfile, wf.Vars)

	tdata, err := os.ReadFile("/tmp/workflow.yaml")
	if err != nil {
		return w, err
	}

	err = yaml.Unmarshal([]byte(tdata), &w)
	if err != nil {
		return w, err
	}

	return w, err
}


func overrideWithEnvs(varsMap map[string]string) error {
	var err error
	for _, v := range os.Environ() {
    split_v := strings.Split(v, "=")
    varsMap[split_v[0]] = split_v[1]
  }

  return err
}

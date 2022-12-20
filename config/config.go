package config

import (
	"fmt"
	"github.com/balmanrawat/cfn-compose/cfn"
	"path/filepath"
	"os"
)

const composeDir string = ".cfn-compose"
const varsTemplate string = "var.yml"
const composeTemplate string = "compose.yml"
var jobCountLimit int = 15
var stackCountLimit int = 30

type ComposeConfig struct {
	Description string            `yaml:"Description"`
	Jobs        map[string]Job    `yaml:"Jobs"`
	Vars        map[string]string `yaml:"Vars"`
}

type Job struct {
	Name        string  `yaml:"Name,omitempty"`
	Description string  `yaml:"Description,omitempty"`
	Stacks      []cfn.Stack `yaml:"Stacks"`
	Order       int     `yaml:"Order"`
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

	if j.Order < 0 || j.Order > 100 {
		return fmt.Errorf("Job Order should be within 0-100 range, found: %d", j.Order)
	}

	// if j.Order < 0 {
	// 	return fmt.Errorf("Job Order Can't be negative value, found: %d", j.Order)
	// }

	for i, stack := range j.Stacks {
		err := stack.Validate(i)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
ComposeConfig is valid when all of the below conditions are true:
- When Job counts is <= jobCoutLimit
- When all jobs are valid
- When all stacks inside the jobs are valid
*/
func (c *ComposeConfig) Validate() error {
	if len(c.Jobs) > jobCountLimit {
		return fmt.Errorf("Job count is %d, should be <= %d", len(c.Jobs), jobCountLimit)
	}

	if len(c.Jobs) <= 0 {
		return fmt.Errorf("Job count is %d, compose config should have at least one job", len(c.Jobs))
	}

	for jname, job := range c.Jobs {
		if err := job.Validate(jname); err != nil {
			return fmt.Errorf("[Job: %s] Error: %s", jname, err.Error())
		}
	}

	return nil
}

func GetComposeConfig(configFile string) (ComposeConfig, error) {
	var cc ComposeConfig
	dir := filepath.Dir(configFile)
	file := filepath.Base(configFile)
	os.Chdir(dir)

	cc, err := parse(file)
	if err != nil {
		fmt.Printf("Failed while fetching compose file: %s\n", err.Error())
	}

	return cc, err
}
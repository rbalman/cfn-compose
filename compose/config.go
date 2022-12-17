package compose

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"github.com/balmanrawat/cfn-compose/cfn"

	"gopkg.in/yaml.v2"
)

type ComposeConfig struct {
	Description string            `yaml:"description"`
	Jobs        map[string]Job    `yaml:"jobs"`
	Vars        map[string]string `yaml:"vars"`
}

var jobCountLimit int = 5
const composeDir string = ".cfn-compose"
const varsTemplate string = "var.yml"
const composeTemplate string = "compose.yml"

type Job struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Stacks      []cfn.Stack `yaml:"stacks"`
	Order       int     `yaml:"order"`
}

var stackCountLimit int = 30
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

//Parsing the configuration file
func Parse(file string) (ComposeConfig, error) {
	if _, err := os.Stat(composeDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(composeDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	var cc ComposeConfig
	data, err := os.ReadFile(file)
	if err != nil {
		return cc, err
	}

	vars, err := prepareVariables(data)
	if err != nil {
		return cc, err
	}

	t, err := template.New("ComposeConfigTemplate").Funcs(template.FuncMap{
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

			return strings.TrimSpace(stdout.String())
		},
	}).Parse(string(data))
	if err != nil {
		return cc, err
	}

	final_template_file, err := os.Create(composeDir + "/" + composeTemplate)
	if err != nil {
		return cc, err
	}
	defer final_template_file.Close()
	t.Execute(final_template_file, vars)

	varsData, err := os.ReadFile(composeDir + "/" + composeTemplate)
	if err != nil {
		return cc, err
	}

	err = yaml.Unmarshal([]byte(varsData), &cc)
	if err != nil {
		return cc, err
	}

	cc.Vars = vars

	return cc, err
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

	vars_file, err := os.Create(composeDir + "/" + varsTemplate)
	if err != nil {
		return varStruct.Vars, err
	}

	defer vars_file.Close()
	err = t.Execute(vars_file, nil)
	if err != nil {
		return varStruct.Vars, err
	}

	varFileData, err := os.ReadFile(composeDir + "/" + varsTemplate)
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

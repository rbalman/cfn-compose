package config

import (
	"fmt"
	"github.com/rbalman/cfn-compose/cfn"
	"os"
	"path/filepath"
)

const composeDir string = ".cfn-compose"
const composeTemplate string = "compose.yml"
<<<<<<< HEAD

=======
>>>>>>> 3d62c37a9f9b502078c157129367e517cb2b4a00
var flowCountLimit int = 50
var stackCountLimit int = 50

type ComposeConfig struct {
	Description string            `yaml:"Description"`
	Flows       map[string]Flow   `yaml:"Flows"`
	Vars        map[string]string `yaml:"Vars"`
}

type Flow struct {
	Name        string      `yaml:"Name,omitempty"`
	Description string      `yaml:"Description,omitempty"`
	Stacks      []cfn.Stack `yaml:"Stacks"`
	Order       int         `yaml:"Order"`
}

/*
Flow is valid when all of the below conditions are true:
- When Stack counts is <= stackCountLimit
- order property should be a valid unsigned integer
- When all stacks are valid
*/
func (j *Flow) Validate(name string) error {
	if len(j.Stacks) >= stackCountLimit || len(j.Stacks) == 0 {
		return fmt.Errorf("Stack count is %d for Flow: %s, should be '> 0 and <= %d'", len(j.Stacks), name, stackCountLimit)
	}

	if j.Order < 0 || j.Order > 100 {
		return fmt.Errorf("Flow Order should be within 0-100 range, found: %d", j.Order)
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
ComposeConfig is valid when all of the below conditions are true:
- When Flow counts is <= flowCoutLimit
- When all flows are valid
- When all stacks inside the flows are valid
*/
func (c *ComposeConfig) Validate() error {
	if len(c.Flows) > flowCountLimit {
		return fmt.Errorf("Flow count is %d, should be <= %d", len(c.Flows), flowCountLimit)
	}

	if len(c.Flows) <= 0 {
		return fmt.Errorf("Flow count is %d, compose config should have at least one flow", len(c.Flows))
	}

	for jname, flow := range c.Flows {
		if err := flow.Validate(jname); err != nil {
			return fmt.Errorf("[Flow: %s] Error: %s", jname, err.Error())
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

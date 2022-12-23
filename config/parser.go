package config

import (
	"errors"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func parse(file string) (ComposeConfig, error) {
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

	vars, err := extractVars(data)
	if err != nil {
		return cc, err
	}

	t, err := template.New("ComposeConfigTemplate").Parse(string(data))
	if err != nil {
		return cc, err
	}

	compose_state_file, err := os.Create(composeDir + "/" + composeTemplate)
	if err != nil {
		return cc, err
	}
	defer compose_state_file.Close()
	t.Execute(compose_state_file, vars)

	composeData, err := os.ReadFile(composeDir + "/" + composeTemplate)
	if err != nil {
		return cc, err
	}

	err = yaml.Unmarshal([]byte(composeData), &cc)
	if err != nil {
		return cc, err
	}

	cc.Vars = vars

	return cc, err
}

func extractVars(data []byte) (map[string]string, error) {
	vars := struct {
		Vmap map[string]string `yaml:"Vars"`
	}{}

	err := yaml.Unmarshal(data, &vars)
	if err != nil {
		return vars.Vmap, err
	}

	err = overrideWithEnvs(vars.Vmap)
	return vars.Vmap, err
}

func overrideWithEnvs(varsMap map[string]string) error {
	var err error
	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		varsMap[split_v[0]] = split_v[1]
	}

	return err
}

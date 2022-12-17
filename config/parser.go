package config

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

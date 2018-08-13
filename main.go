package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
)

var execCommand = exec.Command

func runCredhub(arguments []string) ([]byte, error) {
	cmd := execCommand("credhub", arguments...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	return cmd.CombinedOutput()
}

var credhubListEntry = struct {
	Credentials []struct {
		Name string `json:"name"`
	} `json:"credentials"`
}{}

type entry struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type credentials struct {
	Credentials []entry `json:"credentials"`
}

var entries []entry

func main() {
	app := cli.NewApp()
	app.Name = "credhubbackup"
	app.Usage = "a utility to backup credhub stored credentials"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		flist, _ := runCredhub([]string{"find", "-j"})

		err := json.Unmarshal(flist, &credhubListEntry)
		if err != nil {
			return err
		}

		for _, credential := range credhubListEntry.Credentials {
			values, _ := runCredhub([]string{"get", "-j", "-n", credential.Name})
			var e entry
			err1 := json.Unmarshal(values, &e)
			if err1 != nil {
				return err1
			}

			entries = append(entries, e)
		}
		var creds credentials
		creds.Credentials = entries
		output, err := yaml.Marshal(&creds)
		if err != nil {
			return err
		}
		fmt.Print(string(output))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

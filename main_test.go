package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

/*
	credhub login is required

	{
		credentials: [
		{
			"name": "/concourse/main/pipeline/concourse_password",
			"version_created_at": "2018-07-27T11:10:58Z"
		}
		]
	}
	->
	- name: "/concourse/main/pipeline/concourse_password"
	type: value
	value: fjsdklflsdjklfklsdf
*/

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const credhubRunResult = `{
	credentials: [
	{
		"name": "/concourse/main/pipeline/concourse_password",
		"version_created_at": "2018-07-27T11:10:58Z"
	}
	]
}`

func TestRunCredhub(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	out, err := runCredhub([]string{"find", "-j"})
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
	readableOut := string(out)
	if readableOut != credhubRunResult {
		t.Errorf("Expected %q, got %q", credhubRunResult, readableOut)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, credhubRunResult)
	os.Exit(0)
}

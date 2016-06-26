package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestMain_nodockerfile(t *testing.T) {
	reset, err := cd("./testdata/nodockerfile")
	if err != nil {
		t.Fatal(err)
	}
	defer reset()

	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint"}, &b)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != ExitCodeSuccess {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeSuccess, exitCode)
	}
	if len(b.Bytes()) != 0 {
		t.Errorf("should be no output, but got `%s`", b.String())
	}
}

func TestMain_WithHasSomeProblem(t *testing.T) {
	reset, err := cd("./testdata/someproblem")
	if err != nil {
		t.Fatal(err)
	}
	defer reset()

	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint"}, &b)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != ExitCodeHasProblem {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeHasProblem, exitCode)
	}
	if len(b.Bytes()) == 0 {
		t.Error("should detect some issues, but any issues not exist")
	}
}

func TestMain_WithJSONFormatter(t *testing.T) {
	reset, err := cd("./testdata/someproblem")
	if err != nil {
		t.Fatal(err)
	}
	defer reset()

	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint", "--formatter=json"}, &b)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != ExitCodeHasProblem {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeHasProblem, exitCode)
	}

	ps := []Problem{}
	err = json.NewDecoder(&b).Decode(&ps)
	if err != nil {
		t.Fatal(err)
	}

	if len(ps) == 0 {
		t.Errorf("Should return problems, but got %v", ps)
	}
}

func TestMain_WithNotExistConfig(t *testing.T) {
	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint", "--config=does-not-exist-dflint.yml"}, &b)
	if err == nil {
		t.Error("Error should not be nil, but got nil")
	}
	if exitCode != ExitCodeError {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeError, exitCode)
	}
	if !strings.Contains(err.Error(), "does-not-exist-dflint.yml") {
		t.Errorf("error should be not found config file, but got %s", err.Error())
	}
	if len(b.Bytes()) != 0 {
		t.Errorf("should be no output, but got `%s`", b.String())
	}
}

func TestMain_WithImplicitConfigFile(t *testing.T) {
	reset, err := cd("./testdata/with_dflint_yml")
	if err != nil {
		t.Fatal(err)
	}
	defer reset()

	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint"}, &b)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != ExitCodeHasProblem {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeHasProblem, exitCode)
	}

	if len(b.Bytes()) == 0 {
		t.Error("should detect some issues, but any issues not exist")
	}
	for _, line := range strings.Split(b.String(), "\n") {
		if strings.Contains(line, "DowncaseInstruction") {
			t.Errorf("Type should not be DowncaseInstruction. but got %s", line)
		}
	}
}

func TestMain_WithExplicitConfigFile(t *testing.T) {
	reset, err := cd("./testdata/with_dflint_yml")
	if err != nil {
		t.Fatal(err)
	}
	defer reset()

	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint", "-c", "explicit-dflint.yml"}, &b)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != ExitCodeHasProblem {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeHasProblem, exitCode)
	}

	if len(b.Bytes()) == 0 {
		t.Error("should detect some issues, but any issues not exist")
	}
	for _, line := range strings.Split(b.String(), "\n") {
		if strings.Contains(line, "YesOption") {
			t.Errorf("Type should not be YesOption. but got %s", line)
		}
	}
}

func TestMain_WithUnknownFormatter(t *testing.T) {
	var b bytes.Buffer
	exitCode, err := Main([]string{"dflint", "-f", "hogehoge"}, &b)
	if err == nil {
		t.Fatal("should error, but got nil")
	}
	if exitCode != ExitCodeError {
		t.Errorf("Exit code should be %d, but got %d", ExitCodeError, exitCode)
	}

	if err.Error() != "hogehoge formatter doesn't exist." {
		t.Error(err)
	}
}

func cd(dir string) (func(), error) {
	current, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(dir)
	if err != nil {
		return nil, err
	}

	return func() {
		os.Chdir(current)
	}, nil
}

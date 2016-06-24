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
	err = Main([]string{"dflint"}, &b)
	if err != nil {
		t.Fatal(err)
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
	err = Main([]string{"dflint"}, &b)
	if err != nil {
		t.Fatal(err)
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
	err = Main([]string{"dflint", "--formatter=json"}, &b)
	if err != nil {
		t.Fatal(err)
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
	err := Main([]string{"dflint", "--config=does-not-exist-dflint.yml"}, &b)
	if err == nil {
		t.Error("Error should not be nil, but got nil")
	}
	if !strings.Contains(err.Error(), "does-not-exist-dflint.yml") {
		t.Errorf("error should be not found config file, but got %s", err.Error())
	}
	if len(b.Bytes()) != 0 {
		t.Errorf("should be no output, but got `%s`", b.String())
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

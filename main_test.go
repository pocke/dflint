package main

import (
	"bytes"
	"os"
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

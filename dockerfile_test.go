package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDockerfileNodes(t *testing.T) {
	d, err := newDockerfile([]byte(`FROM busybox
RUN ls
RUN echo 'hoge'
ENV foo bar
RUN sl
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	res := d.Nodes("run")
	if len(res) != 3 {
		t.Errorf("Len should be 3, but got %d", len(res))
	}
}

func TestDockerfileAnalyze(t *testing.T) {
	f, err := newTempDockerfile(`FROM busybox`)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	_, err = Analyze(f.Name(), &Config{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDockerfileAnalyze_WithDisabledRule(t *testing.T) {
	f, err := newTempDockerfile(`FROM busybox
run ls
RUN yum install nginx`)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	ps, err := Analyze(f.Name(), &Config{
		IgnoreRules: []string{"YesOption"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(ps) != 1 {
		t.Errorf("Expected len == 1, but got %d", len(ps))
	}

	for _, p := range ps {
		if p.Type == "YesOption" {
			t.Error("`YesOption` is ignored. but got")
		}
	}
}

func TestDockerfileAnalyze_WithSyntaxError(t *testing.T) {
}

// --- test helper

type tempDockerfile struct {
	*os.File
}

func (f *tempDockerfile) Close() error {
	f.File.Close()
	return os.Remove(f.Name())
}

func newTempDockerfile(value string) (*tempDockerfile, error) {
	f, err := ioutil.TempFile("", "dflint-test")
	if err != nil {
		return nil, err
	}

	f.Write([]byte(`FROM busybox
run ls
RUN yum install nginx`))
	return &tempDockerfile{f}, nil
}

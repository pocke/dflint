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
	f, err := ioutil.TempFile("", "dflint-test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	defer f.Close()

	f.Write([]byte(`FROM busybox`))
	_, err = Analyze(f.Name(), &Config{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDockerfileAnalyze_WithDisabledRule(t *testing.T) {
	f, err := ioutil.TempFile("", "dflint-test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	defer f.Close()

	f.Write([]byte(`FROM busybox
run ls
RUN yum install nginx`))

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

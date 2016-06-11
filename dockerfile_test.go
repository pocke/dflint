package main

import (
	"bytes"
	"testing"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func TestDockerfileNodes(t *testing.T) {
	d := &Dockerfile{
		Content: []byte(`FROM busybox
RUN ls
RUN echo 'hoge'
ENV foo bar
RUN sl
`),
		Path: "/dev/null",
	}

	n, err := parser.Parse(bytes.NewBuffer(d.Content))
	if err != nil {
		t.Fatal(err)
	}
	d.AST = n

	res := d.Nodes("run")
	if len(res) != 3 {
		t.Errorf("Len should be 3, but got %d", len(res))
	}
}

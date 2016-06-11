package main

import (
	"bytes"
	"io/ioutil"

	"github.com/docker/docker/builder/dockerfile/parser"
)

type Dockerfile struct {
	Content []byte
	AST     *parser.Node
}

func NewDockerfile(fpath string) (*Dockerfile, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	n, err := parser.Parse(buf)
	if err != nil {
		return nil, err
	}

	return &Dockerfile{
		Content: b,
		AST:     n,
	}, nil
}

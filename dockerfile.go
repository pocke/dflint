package main

import (
	"bytes"
	"io/ioutil"

	"github.com/docker/docker/builder/dockerfile/parser"
)

type Dockerfile struct {
	Content []byte
	AST     *parser.Node
	Path    string
}

func Analyze(fpath string) ([]Problem, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	n, err := parser.Parse(buf)
	if err != nil {
		return nil, err
	}

	d := &Dockerfile{
		Content: b,
		AST:     n,
		Path:    fpath,
	}

	res := make([]Problem, 0)
	for _, r := range Rules {
		res = append(res, r(d)...)
	}
	return res, nil
}

func (d *Dockerfile) Nodes(instruction string) []*parser.Node {
	res := make([]*parser.Node, 0)
	for _, n := range d.AST.Children {
		if n.Value == instruction {
			res = append(res, n)
		}
	}
	return res
}

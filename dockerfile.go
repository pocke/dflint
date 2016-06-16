package main

import (
	"bytes"
	"io/ioutil"

	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/mvdan/sh"
)

type Dockerfile struct {
	Content []byte
	AST     *parser.Node
	Path    string
}

func newDockerfile(content []byte, path string) (*Dockerfile, error) {
	n, err := parser.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	return &Dockerfile{
		Content: content,
		AST:     n,
		Path:    path,
	}, nil
}

func Analyze(fpath string) ([]Problem, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	d, err := newDockerfile(b, fpath)
	if err != nil {
		return []Problem{{Message: "Syntax Error"}}, nil // TODO
	}

	res := make([]Problem, 0)
	for _, r := range Rules {
		res = append(res, r.Analyze(d)...)
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

func nodeIndex(n []*parser.Node, target *parser.Node) int {
	for i, v := range n {
		if v == target {
			return i
		}
	}
	return -1
}

func parseSh(n *parser.Node) (*sh.File, error) {
	r := bytes.NewReader([]byte(n.Next.Value))
	return sh.Parse(r, "", sh.ParseComments)
}

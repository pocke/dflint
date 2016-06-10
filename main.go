package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func main() {
	err := Main(os.Args)
	if err != nil {
		panic(err)
	}
}

func Main(args []string) error {
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)

	n, err := parser.Parse(buf)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(n)
	// showNode(n)
}

func showNode(n *parser.Node) {
	if n == nil {
		return
	}
	fmt.Println(n.Value)
	for _, child := range n.Children {
		showNode(child)
	}
	showNode(n.Next)
}

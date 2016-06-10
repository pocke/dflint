package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/builder/dockerfile/parser"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := parser.Parse(f)
	if err != nil {
		panic(err)
	}

	showNode(n)
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

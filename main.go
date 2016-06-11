package main

import (
	"encoding/json"
	"os"
)

func main() {
	err := Main(os.Args)
	if err != nil {
		panic(err)
	}
}

func Main(args []string) error {
	d, err := NewDockerfile(args[1])
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(d.AST)
}

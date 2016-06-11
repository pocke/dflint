package main

import (
	"fmt"
	"os"
)

func main() {
	err := Main(os.Args)
	if err != nil {
		panic(err)
	}
}

func Main(args []string) error {
	problems, err := Analyze(args[1])
	if err != nil {
		return err
	}

	for _, p := range problems {
		fmt.Println(p)
	}
	return nil
}

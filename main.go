package main

import "os"

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

	FormatJSON(problems, os.Stdout)
	return nil
}

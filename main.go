package main

import (
	"fmt"
	"os"

	"github.com/ogier/pflag"
)

func main() {
	err := Main(os.Args)
	if err != nil {
		panic(err)
	}
}

func Main(args []string) error {
	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fmtrName := ""
	fs.StringVarP(&fmtrName, "formatter", "f", "json", "Output Formatter")
	fs.Parse(args[1:])

	ps := make([]Problem, 0)
	for _, f := range fs.Args() {
		problems, err := Analyze(f)
		if err != nil {
			return err
		}
		ps = append(ps, problems...)
	}

	fmtr, ok := Formatters[fmtrName]
	if !ok {
		return fmt.Errorf("%s formatter doesn't exist.", fmtrName)
	}
	fmtr(ps, os.Stdout)

	return nil
}

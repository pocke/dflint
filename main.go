package main

import (
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar"
	"github.com/ogier/pflag"
)

func main() {
	err := Main(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fmtrName := ""
	fs.StringVarP(&fmtrName, "formatter", "f", "default", "Output Formatter")
	fs.Parse(args[1:])

	ps := make([]Problem, 0)
	targets := fs.Args()
	if len(targets) == 0 {
		var err error
		targets, err = doublestar.Glob("**/Dockerfile")
		if err != nil {
			return err
		}
	}

	for _, f := range targets {
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

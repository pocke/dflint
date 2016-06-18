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
	fmtrName, targets := ParseArgs(args)
	if len(targets) == 0 {
		var err error
		targets, err = doublestar.Glob("**/Dockerfile")
		if err != nil {
			return err
		}
	}

	c, err := ParseConfig()
	if err != nil {
		return err
	}

	ps := make([]Problem, 0)
	for _, f := range targets {
		problems, err := Analyze(f, c)
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

func ParseArgs(args []string) (fmtrName string, arguments []string) {
	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fs.StringVarP(&fmtrName, "formatter", "f", "default", "Specify output formatter. []")
	fs.Parse(args[1:])

	return fmtrName, fs.Args()
}

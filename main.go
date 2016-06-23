package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/ogier/pflag"
)

const DEFAULT_CONF_PATH = "./.dflint.yaml"

func main() {
	err := Main(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	fmtrName, confPath, targets, err := ParseArgs(args)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		targets, err = doublestar.Glob("**/Dockerfile")
		if err != nil {
			return err
		}
	}

	c, err := ParseConfig(confPath)
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

func ParseArgs(args []string) (fmtrName, confPath string, arguments []string, err error) {
	fmtrNames := []string{}
	for name := range Formatters {
		fmtrNames = append(fmtrNames, name)
	}

	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fs.StringVarP(&fmtrName, "formatter", "f", "default", fmt.Sprintf("Specify output formatter. [%s]", strings.Join(fmtrNames, ", ")))
	fs.StringVarP(&confPath, "config", "c", "", "Path of Configuration file (default \"./.dflint.yml\")")
	fs.Parse(args[1:])

	if confPath == "" {
		confPath = DEFAULT_CONF_PATH
	} else {
		// Check confPath exists when specified.
		if !fileExists(confPath) {
			return "", "", nil, fmt.Errorf("%s doesn't exist", confPath)
		}
	}

	arguments = fs.Args()
	return
}

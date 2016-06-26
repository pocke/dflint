package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bmatcuk/doublestar"
	"github.com/ogier/pflag"
)

const (
	ExitCodeSuccess    = 0
	ExitCodeHasProblem = 1
	ExitCodeError      = 2
)

const DEFAULT_CONF_PATH = "./.dflint.yml"

func main() {
	exitCode, err := Main(os.Args, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

// Main returns an exit code and an error.
func Main(args []string, out io.Writer) (int, error) {
	cmdArg, err := ParseArgs(args)
	if err != nil {
		return ExitCodeError, err
	}
	targets := cmdArg.Arguments

	if len(targets) == 0 {
		targets, err = doublestar.Glob("**/Dockerfile")
		if err != nil {
			return ExitCodeError, err
		}
	}

	c, err := ParseConfig(cmdArg.ConfigPath)
	if err != nil {
		return ExitCodeError, err
	}

	ps := make([]Problem, 0)
	for _, f := range targets {
		problems, err := Analyze(f, c)
		if err != nil {
			return ExitCodeError, err
		}
		ps = append(ps, problems...)
	}

	cmdArg.Formatter(ps, out)

	if len(ps) == 0 {
		return ExitCodeSuccess, nil
	} else {
		return ExitCodeHasProblem, nil
	}
}

type CmdArg struct {
	Formatter  FormatFunc
	ConfigPath string
	Arguments  []string
}

func ParseArgs(args []string) (*CmdArg, error) {
	res := new(CmdArg)
	var fmtrName string

	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fs.StringVarP(&fmtrName, "formatter", "f", "default", fmt.Sprintf("Specify output formatter. [%s]", FormatterNames()))
	fs.StringVarP(&res.ConfigPath, "config", "c", "", "Path of Configuration file (default \"./.dflint.yml\")")
	fs.Parse(args[1:])

	if res.ConfigPath == "" {
		res.ConfigPath = DEFAULT_CONF_PATH
	} else {
		// Check confPath exists when specified.
		if !fileExists(res.ConfigPath) {
			return nil, fmt.Errorf("%s doesn't exist", res.ConfigPath)
		}
	}

	fmtr, ok := Formatters[fmtrName]
	if !ok {
		return nil, fmt.Errorf("%s formatter doesn't exist.", fmtrName)
	}
	res.Formatter = fmtr

	res.Arguments = fs.Args()
	return res, nil
}

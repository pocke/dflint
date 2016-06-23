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
	cmdArg, err := ParseArgs(args)
	if err != nil {
		return err
	}
	targets := cmdArg.Arguments

	if len(targets) == 0 {
		targets, err = doublestar.Glob("**/Dockerfile")
		if err != nil {
			return err
		}
	}

	c, err := ParseConfig(cmdArg.ConfigPath)
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

	fmtr, ok := Formatters[cmdArg.FormatterName]
	if !ok {
		return fmt.Errorf("%s formatter doesn't exist.", cmdArg.FormatterName)
	}
	fmtr(ps, os.Stdout)

	return nil
}

type CmdArg struct {
	FormatterName string
	ConfigPath    string
	Arguments     []string
}

func ParseArgs(args []string) (*CmdArg, error) {
	res := new(CmdArg)

	fmtrNames := []string{}
	for name := range Formatters {
		fmtrNames = append(fmtrNames, name)
	}

	fs := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	fs.StringVarP(&res.FormatterName, "formatter", "f", "default", fmt.Sprintf("Specify output formatter. [%s]", strings.Join(fmtrNames, ", ")))
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

	res.Arguments = fs.Args()
	return res, nil
}

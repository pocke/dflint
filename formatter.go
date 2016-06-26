package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type FormatFunc func([]Problem, io.Writer)

func FormatDefault(ps []Problem, w io.Writer) {
	for _, p := range ps {
		fmt.Fprintf(w, "%s:%d:%d: [%s] %s\n", p.Path, p.Line, p.Column, p.Type, p.Message)
	}
}

func FormatJSON(ps []Problem, w io.Writer) {
	json.NewEncoder(w).Encode(ps)
}

var Formatters = map[string]FormatFunc{
	"default": FormatDefault,
	"json":    FormatJSON,
}

func FormatterNames() string {
	fmtrNames := []string{}
	for name := range Formatters {
		fmtrNames = append(fmtrNames, name)
	}
	return strings.Join(fmtrNames, ", ")
}

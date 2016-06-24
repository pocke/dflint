package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func FormatDefault(ps []Problem, w io.Writer) {
	for _, p := range ps {
		fmt.Fprintf(w, "%s:%d:%d: [%s] %s\n", p.Path, p.Line, p.Column, p.Type, p.Message)
	}
}

func FormatJSON(ps []Problem, w io.Writer) {
	json.NewEncoder(w).Encode(ps)
}

var Formatters = map[string]func([]Problem, io.Writer){
	"default": FormatDefault,
	"json":    FormatJSON,
}

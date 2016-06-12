package main

import (
	"encoding/json"
	"io"
)

func FormatJSON(ps []Problem, w io.Writer) {
	json.NewEncoder(w).Encode(ps)
}

var Formatters = map[string]func([]Problem, io.Writer){
	"json": FormatJSON,
}

package main

import (
	"regexp"
	"strings"
)

var Rules = []func(*Dockerfile) []Problem{
	// Instruction should be upcase.
	func(d *Dockerfile) []Problem {
		res := make([]Problem, 0)

		re := regexp.MustCompile(`^\s*\w+`)
		for _, n := range d.AST.Children {
			// n.Value is downcase always. So, should compare n.Original
			ins := re.FindString(n.Original)
			if strings.ToUpper(ins) != ins {
				res = append(res, Problem{
					Line:    n.StartLine,
					Column:  0, // TODO
					Length:  len(n.Value),
					Path:    d.Path,
					Type:    "DowncaseInstruction",
					Message: "Instruction should be upcase.",
				})
			}
		}

		return res
	},
}

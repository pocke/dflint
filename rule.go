package main

import (
	"regexp"
	"strings"
)

type Rule struct {
	Type string
	f    func(*Rule, *Dockerfile) []Problem
}

func (r *Rule) MakeProblem(line, col, len int, dockerfile *Dockerfile, msg string) Problem {
	return Problem{
		Line:    line,
		Column:  col,
		Length:  len,
		Path:    dockerfile.Path,
		Type:    r.Type,
		Message: msg,
	}
}

func (r *Rule) Analyze(d *Dockerfile) []Problem {
	return r.f(r, d)
}

var Rules = []Rule{
	{
		Type: "DowncaseInstruction",
		f: func(r *Rule, d *Dockerfile) []Problem {
			res := make([]Problem, 0)

			re := regexp.MustCompile(`^\s*\w+`)
			for _, n := range d.AST.Children {
				// n.Value is downcase always. So, should compare n.Original
				ins := re.FindString(n.Original)
				if strings.ToUpper(ins) != ins {
					res = append(res, r.MakeProblem(
						n.StartLine,
						0, // TODO
						len(n.Value),
						d,
						"Instruction should be upcase.",
					))
				}
			}

			return res
		},
	},
}

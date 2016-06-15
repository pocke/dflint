package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mvdan/sh"
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

	{
		Type: "FROMShouldBeFirst",
		f: func(r *Rule, d *Dockerfile) []Problem {
			froms := d.Nodes("from")

			if len(froms) == 0 {
				return []Problem{r.MakeProblem(
					0,
					0,
					0,
					d,
					"FROM Instruction doesn't found.",
				)}
			}

			if len(froms) == 1 && nodeIndex(d.AST.Children, froms[0]) == 0 {
				return []Problem{}
			}

			if len(froms) > 1 {
				res := make([]Problem, 0, len(froms))
				for _, n := range froms {
					res = append(res, r.MakeProblem(
						n.StartLine,
						0, // TODO
						0, // TODO
						d,
						"Too many FROM Instruction",
					))

				}
				return res
			}

			return []Problem{r.MakeProblem(
				froms[0].StartLine,
				0, // TODO
				0, // TODO
				d,
				"FROM Instruction should be at first line",
			)}
		},
	},

	{
		Type: "ExportInRUN",
		f: func(r *Rule, d *Dockerfile) []Problem {
			res := make([]Problem, 0)

			for _, v := range d.Nodes("run") {
				f, err := parseSh(v)
				if err != nil {
					// the err is syntax error of shell. so, the err shoulde be handled in other rule.
					continue
				}

				w := ShWalker{
					onCallExpr: func(s sh.CallExpr) {
						if len(s.Args) != 2 {
							return
						}
						if l, ok := s.Args[0].Parts[0].(sh.Lit); !ok || l.Value != "export" {
							return
						}
						env, ok := s.Args[1].Parts[0].(sh.Lit)
						re := regexp.MustCompile(`^(\w+)\=(.+)$`)
						if !ok || !re.MatchString(env.Value) {
							return
						}

						// TODO: line,col,length...
						res = append(res, r.MakeProblem(
							v.StartLine,
							0,
							0,
							d,
							fmt.Sprintf("Does not work `export` in RUN instruction. Use `ENV %s` instead of this.", env.Value),
						))
					},
				}

				for _, stmt := range f.Stmts {
					w.Walk(stmt)
				}
			}

			return res
		},
	},
}

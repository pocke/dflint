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
					onCallExpr: func(s *sh.CallExpr) {
						if len(s.Args) != 2 {
							return
						}

						if !callExprEq(s, 0, "export") {
							return
						}
						env, ok := s.Args[1].Parts[0].(*sh.Lit)
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

	{
		Type: "ShellSyntaxError",
		f: func(r *Rule, d *Dockerfile) []Problem {
			res := make([]Problem, 0)

			for _, v := range d.Nodes("run") {
				_, err := parseSh(v)
				if err != nil {
					// TODO: line, col, etc...
					res = append(res, r.MakeProblem(
						v.StartLine,
						0,
						0,
						d,
						fmt.Sprintf("Shell Syntax Error: %s", err.Error()),
					))
				}
			}

			return res
		},
	},

	// TODO: Support apt-get, etc
	{
		Type: "YesOption",
		f: func(r *Rule, d *Dockerfile) []Problem {
			res := make([]Problem, 0)

			for _, v := range d.Nodes("run") {
				f, err := parseSh(v)
				if err != nil {
					continue
				}

				w := ShWalker{
					onCallExpr: func(s *sh.CallExpr) {
						if !(callExprEq(s, 0, "yum") || (callExprEq(s, 0, "sudo") && callExprEq(s, 0, "yum"))) {
							return
						}
						for idx := range s.Args {
							if callExprEq(s, idx, "-y") {
								return
							}
						}
						// TODO: line, col, ...
						res = append(res, r.MakeProblem(
							v.StartLine,
							0,
							0,
							d,
							"`-y` option is required with yum.",
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

	{
		Type: "UnknownInstruction",
		f: func(r *Rule, d *Dockerfile) []Problem {
			AVAILABLE_INSTS := []string{
				"FROM",
				"MAINTAINER",
				"RUN",
				"CMD",
				"LABEL",
				"EXPOSE",
				"ENV",
				"ADD",
				"COPY",
				"ENTRYPOINT",
				"VOLUME",
				"USER",
				"WORKDIR",
				"ARG",
				"ONBUILD",
				"STOPSIGNAL",
				"HEALTHCHECK",
				"SHELL",
			}
			res := make([]Problem, 0)

			for _, n := range d.AST.Children {
				inst := strings.ToUpper(n.Value)
				if hasString(AVAILABLE_INSTS, inst) {
					continue
				}

				res = append(res, r.MakeProblem(
					n.StartLine,
					0, // TODO
					len(n.Value),
					d,
					fmt.Sprintf("%s is unknown instruction", inst),
				))
			}

			return res
		},
	},
}

func callExprEq(s *sh.CallExpr, idx int, target string) bool {
	if len(s.Args) < idx+1 {
		return false
	}
	l, ok := s.Args[idx].Parts[0].(*sh.Lit)
	if !ok {
		return false
	}
	return l.Value == target
}

func hasString(slice []string, t string) bool {
	for _, v := range slice {
		if v == t {
			return true
		}
	}
	return false
}

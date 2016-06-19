package main

import (
	"fmt"
	"testing"
)

func TestRuleDowncaseInstruction(t *testing.T) {
	ps := analyzeOneRule("DowncaseInstruction", `FROM busybox
run ls
`)
	if len(ps) != 1 {
		t.Errorf("Problems should be 1, but got %d", len(ps))
	}
	if ps[0].Message != "Instruction should be upcase." {
		t.Errorf("Message should be 'Instruction should be upcase.', but got %s", ps[0].Message)
	}
	if ps[0].Type != "DowncaseInstruction" {
		t.Errorf("Type should be 'DowncaseInstruction', but got %s", ps[0].Type)
	}
}

func TestRuleFROMShouldBeFirst_valid(t *testing.T) {
	ps := analyzeOneRule("FROMShouldBeFirst", `FROM busybox
RUN ls
`)
	if len(ps) != 0 {
		t.Errorf("should find no problems, but got %v", ps)
	}
}

func TestRuleFROMShouldBeFirst_doesntHaveFROM(t *testing.T) {
	ps := analyzeOneRule("FROMShouldBeFirst", `RUN ls`)
	if len(ps) != 1 {
		t.Errorf("should find one problem, but got %v", ps)
	}
	if ps[0].Type != "FROMShouldBeFirst" {
		t.Errorf("Type should be 'FROMShouldBeFirst', but got %s", ps[0].Type)
	}
}

func TestRuleFROMShouldBeFirst_hasManyFROM(t *testing.T) {
	ps := analyzeOneRule("FROMShouldBeFirst", `FROM scratch
FROM busybox
`)
	if len(ps) != 2 {
		t.Errorf("should find two problems, but got %v", ps)
	}
	if ps[0].Type != "FROMShouldBeFirst" {
		t.Errorf("Type should be 'FROMShouldBeFirst', but got %s", ps[0].Type)
	}
}

func TestRuleFROMShouldBeFirst_FROMIsNotFirst(t *testing.T) {
	ps := analyzeOneRule("FROMShouldBeFirst", `RUN ls
FROM busybox
`)
	if len(ps) != 1 {
		t.Errorf("should find one problem, but got %v", ps)
	}
	if ps[0].Type != "FROMShouldBeFirst" {
		t.Errorf("Type should be 'FROMShouldBeFirst', but got %s", ps[0].Type)
	}
	if ps[0].Line != 2 {
		t.Errorf("Line should be 2, but got %d", ps[0].Line)
	}
}

func TestRuleExportInRUN(t *testing.T) {
	ps := analyzeOneRule("ExportInRUN", `FROM busybox
RUN export FOO=BAR
`)
	if len(ps) != 1 {
		t.Errorf("should find one problem, but got %v", ps)
	}
	if ps[0].Type != "ExportInRUN" {
		t.Errorf("Type should be 'ExportInRUN', but got %s", ps[0].Type)
	}
	if ps[0].Line != 2 {
		t.Errorf("Line should be 2, but got %d", ps[0].Line)
	}
}

func TestRuleShellSyntaxError(t *testing.T) {
	ps := analyzeOneRule("ShellSyntaxError", `FROM busybox
RUN echo hoge
RUN foo &&
`)
	if len(ps) != 1 {
		t.Errorf("should find one problem, but got %v", ps)
	}
	if ps[0].Type != "ShellSyntaxError" {
		t.Errorf("Type should be 'ShellSyntaxError', but got %s", ps[0].Type)
	}
	if ps[0].Line != 3 {
		t.Errorf("Line should be 3, but got %d", ps[0].Line)
	}
}

func TestRuleYesOption(t *testing.T) {
	ps := analyzeOneRule("YesOption", `FROM busybox
RUN yum install nginx -y
RUN yum install nginx
RUN echo yum install nginx`)

	if len(ps) != 1 {
		t.Errorf("should find one problems, but got %v", ps)
	}
	if ps[0].Type != "YesOption" {
		t.Errorf("Type should be 'YesOption', but got %s", ps[0].Type)
	}
	if ps[0].Line != 3 {
		t.Errorf("Line should be 3, but got %d", ps[0].Line)
	}

}

func analyzeOneRule(rule, dockerfile string) []Problem {
	r := findRule(rule)
	d, err := newDockerfile([]byte(dockerfile), "/dev/null")
	if err != nil {
		panic(err)
	}

	return r.Analyze(d)
}

func findRule(t string) Rule {
	for _, r := range Rules {
		if r.Type == t {
			return r
		}
	}
	panic(fmt.Sprintf("%s doesn't found", t))
}

package main

import (
	"fmt"
	"testing"
)

func TestRuleDowncaseInstruction(t *testing.T) {
	r := findRule("DowncaseInstruction")
	d, err := newDockerfile([]byte(`FROM busybox
run ls
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
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
	r := findRule("FROMShouldBeFirst")
	d, err := newDockerfile([]byte(`FROM busybox
RUN ls
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
	if len(ps) != 0 {
		t.Errorf("should find no problems, but got %v", ps)
	}
}

func TestRuleFROMShouldBeFirst_doesntHaveFROM(t *testing.T) {
	r := findRule("FROMShouldBeFirst")
	d, err := newDockerfile([]byte(`RUN ls`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
	if len(ps) != 1 {
		t.Errorf("should find one problem, but got %v", ps)
	}
	if ps[0].Type != "FROMShouldBeFirst" {
		t.Errorf("Type should be 'FROMShouldBeFirst', but got %s", ps[0].Type)
	}
}

func TestRuleFROMShouldBeFirst_hasManyFROM(t *testing.T) {
	r := findRule("FROMShouldBeFirst")
	d, err := newDockerfile([]byte(`FROM scratch
FROM busybox
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
	if len(ps) != 2 {
		t.Errorf("should find two problems, but got %v", ps)
	}
	if ps[0].Type != "FROMShouldBeFirst" {
		t.Errorf("Type should be 'FROMShouldBeFirst', but got %s", ps[0].Type)
	}
}

func TestRuleFROMShouldBeFirst_FROMIsNotFirst(t *testing.T) {
	r := findRule("FROMShouldBeFirst")
	d, err := newDockerfile([]byte(`RUN ls
FROM busybox
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
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

func findRule(t string) Rule {
	for _, r := range Rules {
		if r.Type == t {
			return r
		}
	}
	panic(fmt.Sprintf("%s doesn't found", t))
}

func TestRuleExportInRUN(t *testing.T) {
	r := findRule("ExportInRUN")
	d, err := newDockerfile([]byte(`FROM busybox
RUN export FOO=BAR
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
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
	r := findRule("ShellSyntaxError")
	d, err := newDockerfile([]byte(`FROM busybox
RUN echo hoge
RUN foo &&
`), "/dev/null")
	if err != nil {
		t.Error(err)
	}

	ps := r.Analyze(d)
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

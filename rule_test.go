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

func findRule(t string) Rule {
	for _, r := range Rules {
		if r.Type == t {
			return r
		}
	}
	panic(fmt.Sprintf("%s doesn't found", t))
}

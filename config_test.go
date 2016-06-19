package main

import "testing"

func TestConfigIsEnabledRule(t *testing.T) {
	c := &Config{
		IgnoreRules: []string{"YesOption"},
	}

	if !c.IsEnabledRule("ShellSyntaxError") {
		t.Errorf("%s should be enabled, but got false", "ShellSyntaxError")
	}

	if c.IsEnabledRule("YesOption") {
		t.Errorf("%s should be disabled, but got true", "YesOption")
	}
}

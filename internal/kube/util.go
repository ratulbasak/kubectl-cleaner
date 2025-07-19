package kube

import (
	"gopkg.in/yaml.v3"
	"io"
)

func DefaultAnalyzerRules() AnalyzerRules {
	return AnalyzerRules{}
}

func LoadAnalyzerRules(r io.Reader) AnalyzerRules {
	var rules AnalyzerRules
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&rules); err != nil {
		return AnalyzerRules{}
	}
	return rules
}

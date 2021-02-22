package helper

import (
	"github.com/accurics/terrascan/pkg/results"
)

type violations []*results.Violation

func (v violations) Len() int {
	return len(v)
}

func (v violations) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v violations) Less(i, j int) bool {
	if v[i].File < v[j].File {
		return true
	}
	if v[i].File > v[j].File {
		return false
	}

	if v[i].ResourceType < v[j].ResourceType {
		return true
	}

	if v[i].ResourceType > v[j].ResourceType {
		return false
	}

	if v[i].RuleName < v[j].RuleName {
		return true
	}

	if v[i].RuleName > v[j].RuleName {
		return false
	}

	if v[i].ResourceName < v[j].ResourceName {
		return true
	}

	if v[i].ResourceName > v[j].ResourceName {
		return false
	}

	if v[i].LineNumber < v[j].LineNumber {
		return true
	}

	return v[i].LineNumber > v[j].LineNumber
}

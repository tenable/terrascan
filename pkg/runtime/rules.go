package runtime

import (
	"github.com/tenable/terrascan/pkg/config"
)

// read the config file and update scan and skip rules
func (e *Executor) loadRuleSetFromConfig() error {

	// append scan rules
	if len(config.GetScanRules()) > 0 {
		e.scanRules = append(e.scanRules, config.GetScanRules()...)
	}

	// append skip rules
	if len(config.GetSkipRules()) > 0 {
		e.skipRules = append(e.skipRules, config.GetSkipRules()...)
	}

	e.categories = config.GetCategoryList()

	// specify severity of violations to be reported
	if len(config.GetSeverityLevel()) > 0 {
		e.severity = config.GetSeverityLevel()
	}

	return nil
}

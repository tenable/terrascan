package runtime

import (
	"github.com/accurics/terrascan/pkg/config"
	"go.uber.org/zap"
)

// read the config file and update scan and skip rules
func (e *Executor) initRules() error {
	if e.configFile == "" {
		return nil
	}

	configReader, err := config.NewTerrascanConfigReader(e.configFile)
	if err != nil {
		zap.S().Error("error loading config file", zap.Error(err))
		return err
	}

	// append scan rules
	if len(configReader.GetRules().ScanRules) > 0 {
		e.scanRules = append(e.scanRules, configReader.GetRules().ScanRules...)
	}

	// append skip rules
	if len(configReader.GetRules().SkipRules) > 0 {
		e.skipRules = append(e.skipRules, configReader.GetRules().SkipRules...)
	}
	return nil
}

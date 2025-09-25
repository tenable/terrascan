package commons

import (
	"fmt"
	"path/filepath"
)

/*
ErrFmtInvalidModuleConfigs defines error when failed to build unified config
*/
const ErrFmtInvalidModuleConfigs = `failed to build unified config. errors:
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
`

/*
ErrFmtTerraformLoad defines error occurred while loading terraform config dir
*/
const ErrFmtTerraformLoad = `diagnostic errors while loading terraform config dir '%s'. error from terraform:
%s:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
%s:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.
`

/*
ErrStrBuildingUnifiedConfig defines error when failed to build unified config
*/
const ErrStrBuildingUnifiedConfig = `failed to build unified config. errors:
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
`

/*
ErrStrModuleSourceInvalid defines error when failed to build unified config in case of invalid module source
*/
const ErrStrModuleSourceInvalid = `failed to build unified config. errors:
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module cloudfront
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module m1
`

/*
ErrMsgFailedLoadingConfigFile defines error when fails to load the config file
*/
const ErrMsgFailedLoadingConfigFile = `error occurred while loading config file 'not-there'. error: <nil>: Failed to read file; The file "not-there" could not be read.`

/*
ErrMsgFailedLoadingIACFile defines error when fails to load the IAC file
*/
const ErrMsgFailedLoadingIACFile = `failed to load iac file '%s'. error: %s:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
%s:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.`

/*
GenerateInvalidModuleConfigError returns invalid module config error with the given test data dir
*/
func GenerateInvalidModuleConfigError(testDataDir string) error {
	return fmt.Errorf(ErrFmtInvalidModuleConfigs, filepath.Join(testDataDir, "invalid-moduleconfigs", "cloudfront", "sub-cloudfront")) //lint:ignore ST1005 the newline in the error message has been added intentionally
}

/*
GenerateTerraformLoadError returns terraform load error with the given test data dir
*/
func GenerateTerraformLoadError(testDataDir, emptyTfFilePath1 string, emptyTfFilePath2 string) error {
	return fmt.Errorf(ErrFmtTerraformLoad, testDataDir, emptyTfFilePath1, emptyTfFilePath2) //lint:ignore ST1005 the newline in the error message has been added intentionally
}

/*
GenerateErrBuildingUnifiedConfig returns error when fails building the unified config
*/
func GenerateErrBuildingUnifiedConfig(testDataDir string) error {
	return fmt.Errorf(ErrStrBuildingUnifiedConfig, filepath.Join(testDataDir, "depends_on", "live", "log"), filepath.Join(testDataDir, "depends_on", "live", "security")) //lint:ignore ST1005 the newline in the error message has been added intentionally
}

/*
GenerateErrStringModuleSourceInvalid returns error when fails building the unified config due to the invalid source
*/
func GenerateErrStringModuleSourceInvalid(testDataDir string) error {
	return fmt.Errorf(ErrStrModuleSourceInvalid, filepath.Join(testDataDir, "invalid-module-source"), filepath.Join(testDataDir, "invalid-module-source")) //lint:ignore ST1005 the newline in the error message has been added intentionally
}

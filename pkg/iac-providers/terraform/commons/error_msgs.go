package commons

import (
	"fmt"
	"path/filepath"
)

const errFmtInvalidModuleConfigs = `failed to build unified config. errors:
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
`
const errFmtTerraformLoad = `diagnostic errors while loading terraform config dir '%s'. error from terraform:
%s:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
%s:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.
`
const errStringDependsOn = `failed to build unified config. errors:
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
`

const errStrModuleSourceInvalid = `failed to build unified config. errors:
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module cloudfront
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module m1
`

func GenerateInvalidModuleConfigError(testDataDir string) error {
	return fmt.Errorf(errFmtInvalidModuleConfigs, filepath.Join(testDataDir, "invalid-moduleconfigs", "cloudfront", "sub-cloudfront"))
}

func GenerateTerraformLoadError(testDataDir, emptyTfFilePath1 string, emptyTfFilePath2 string) error {
	return fmt.Errorf(errFmtTerraformLoad, testDataDir, emptyTfFilePath1, emptyTfFilePath2)
}

func GenerateErrStringDependsOn(testDataDir string) error {
	return fmt.Errorf(errStringDependsOn, filepath.Join(testDataDir, "depends_on", "live", "log"), filepath.Join(testDataDir, "depends_on", "live", "security"))
}

func GenerateErrStringModuleSourceInvalid(testDataDir string) error {
	return fmt.Errorf(errStrModuleSourceInvalid, filepath.Join(testDataDir, "invalid-module-source"), filepath.Join(testDataDir, "invalid-module-source"))
}

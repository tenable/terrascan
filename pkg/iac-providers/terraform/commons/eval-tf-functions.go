package commons

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	modulePath         = "${path.module}"
	rootModulePath     = "${path.root}"
	currentWorkingPath = "${path.cwd}"
)

func evaluateTemplatefile(exprValue, modfiledir string) (string, error) {
	params := getTemplatefileParams(exprValue)
	interpretedPath, err := interpretFilesystemInfo(params[0], modfiledir)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(interpretedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	templateInfo := string(data)

	for i := 1; i < len(params); i++ {
		key := fmt.Sprintf("${%s}", strings.Split(params[i], "=")[0])
		value := strings.Split(params[i], "=")[1]

		templateInfo = strings.ReplaceAll(templateInfo, key, value)
	}

	return templateInfo, nil
}

func getTemplatefileParams(exprValue string) []string {
	re := regexp.MustCompile(`\((.*)\)`)
	paramString := re.FindString(exprValue)
	paramString = strings.TrimLeft(paramString, "(")
	paramString = strings.TrimRight(paramString, ")")
	paramString = strings.ReplaceAll(paramString, "\"", "")
	paramString = strings.ReplaceAll(paramString, " ", "")
	paramOne := strings.Split(paramString, ",")[0]

	re = regexp.MustCompile(`(,{)(.*)(})`)
	paramTwo := re.FindString(paramString)
	paramTwo = strings.TrimLeft(paramTwo, ",{")
	paramTwo = strings.TrimRight(paramTwo, "}")

	var params []string
	params = append(params, paramOne)
	params = append(params, strings.Split(paramTwo, ",")...)

	return params
}

func interpretFilesystemInfo(fsinfo, modfiledir string) (string, error) {
	if strings.HasPrefix(fsinfo, modulePath) {
		return strings.Replace(fsinfo, modulePath, modfiledir, 1), nil
	}

	if strings.HasPrefix(fsinfo, rootModulePath) {
		return strings.Replace(fsinfo, rootModulePath, modfiledir, 1), nil
	}

	if strings.HasPrefix(fsinfo, currentWorkingPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working dir: %w", err)
		}
		return strings.Replace(fsinfo, modulePath, cwd, 1), nil
	}

	return fsinfo, nil
}

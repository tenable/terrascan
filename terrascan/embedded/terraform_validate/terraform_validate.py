import hcl
import os
import re
import traceback
import json
import logging
import sys

class TerraformPropertyList:

    def __init__(self, validator):
        self.properties = []
        self.validator = validator

    def tfproperties(self):
        return self.properties

    def property(self, property_name):
        propList = TerraformPropertyList(self.validator)
        for prop in self.properties:
            pvList = []
            if type(prop.property_value) is list:
                pvList = prop.property_value
            else:
                pvList.append(prop.property_value)

            wasFound = False
            for pv in pvList:
                if type(pv) is dict and property_name in pv.keys():
                    wasFound = True
                    propList.properties.append(TerraformProperty(prop.resource_type,
                                               "{0}.{1}".format(prop.resource_name, prop.property_name),
                                               property_name,
                                               pv[property_name],
                                               prop.moduleName,
                                               prop.fileName))
                    
            if not wasFound and self.validator.raise_error_if_property_missing:                    
                self.validator.preprocessor.add_failure(
                    "[{0}.{1}] should have property: '{2}'".format(prop.resource_type, "{0}.{1}".format(prop.resource_name, prop.property_name), property_name),
                    prop.moduleName,
                    prop.fileName,
                    self.validator.severity,
                    self.validator.isRuleOverridden,
                    self.validator.overrides,
                    prop.resource_name)

        return propList

    def should_equal_case_insensitive(self, expected_value):
        self.should_equal(expected_value, True)

    def should_equal(self, expected_value, caseInsensitive=False):
        for prop in self.properties:

            expected_value = self.int2str(expected_value)
            prop.property_value = self.int2str(prop.property_value)
            expected_value = self.bool2str(expected_value)
            prop.property_value = self.bool2str(prop.property_value)

            if caseInsensitive:
                # make both actual and expected lower case so case won't matter
                pv = prop.property_value.lower()
                ev = expected_value.lower()
            else:
                pv = prop.property_value
                ev = expected_value

            if pv != ev:
                self.validator.preprocessor.add_failure("[{0}.{1}.{2}] should be '{3}'. Is: '{4}'".format(prop.resource_type,
                                                                                                          prop.resource_name,
                                                                                                          prop.property_name,
                                                                                                          expected_value,
                                                                                                          prop.property_value),
                                                        prop.moduleName,
                                                        prop.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        prop.resource_name)

    def should_not_equal_case_insensitive(self, expected_value):
        self.should_not_equal(expected_value, True)

    def should_not_equal(self, expected_value, caseInsensitive=False):
        for prop in self.properties:

            prop.property_value = self.int2str(prop.property_value)
            expected_value = self.int2str(expected_value)
            expected_value = self.bool2str(expected_value)
            prop.property_value = self.bool2str(prop.property_value)

            if caseInsensitive:
                # make both actual and expected lower case so case won't matter
                pv = prop.property_value.lower()
                ev = expected_value.lower()
            else:
                pv = prop.property_value
                ev = expected_value

            if pv == ev:
                self.validator.preprocessor.add_failure("[{0}.{1}.{2}] should not be '{3}'. Is: '{4}'".format(prop.resource_type,
                                                                                                              prop.resource_name,
                                                                                                              prop.property_name,
                                                                                                              expected_value,
                                                                                                              prop.property_value),
                                                        prop.moduleName,
                                                        prop.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        prop.resource_name)

    def list_should_contain_any(self, values_list):
        if type(values_list) is not list:
            values_list = [values_list]

        for prop in self.properties:
            property_value = prop.property_value
            if type(property_value) is not list:
                property_value = [property_value]
            for pv in property_value:
                if pv not in values_list:
                    if type(prop.property_value) is list:
                        prop.property_value = [str(x) for x in prop.property_value]  # fix 2.6/7
                    self.validator.preprocessor.add_failure("[{0}.{1}.{2}] '{3}' should have been one of '{4}'.".format(prop.resource_type,
                                                                                                                        prop.resource_name,
                                                                                                                        prop.property_name,
                                                                                                                        prop.property_value,
                                                                                                                        values_list),
                                                            prop.moduleName,
                                                            prop.fileName,
                                                            self.validator.severity,
                                                            self.validator.isRuleOverridden,
                                                            self.validator.overrides,
                                                            prop.resource_name)
                    break;

    def list_should_contain(self, values_list):
        if type(values_list) is not list:
            values_list = [values_list]

        for prop in self.properties:

            values_missing = []
            for value in values_list:
                if value not in prop.property_value:
                    values_missing.append(value)

            if len(values_missing) != 0:
                if type(prop.property_value) is list:
                    prop.property_value = [str(x) for x in prop.property_value]  # fix 2.6/7
                self.validator.preprocessor.add_failure("[{0}.{1}.{2}] '{3}' should contain '{4}'.".format(prop.resource_type,
                                                                                                           prop.resource_name,
                                                                                                           prop.property_name,
                                                                                                           prop.property_value,
                                                                                                           values_missing),
                                                        prop.moduleName,
                                                        prop.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        prop.resource_name)

    def list_should_not_contain(self, values_list):
        if type(values_list) is not list:
            values_list = [values_list]

        for prop in self.properties:

            values_missing = []
            for value in values_list:
                if value in prop.property_value:
                    values_missing.append(value)

            if len(values_missing) != 0:
                if type(prop.property_value) is list:
                    prop.property_value = [str(x) for x in prop.property_value]  # fix 2.6/7
                self.validator.preprocessor.add_failure("[{0}.{1}.{2}] '{3}' should not contain '{4}'.".format(prop.resource_type,
                                                                                                               prop.resource_name,
                                                                                                               prop.property_name,
                                                                                                               prop.property_value,
                                                                                                               values_missing),
                                                        prop.moduleName,
                                                        prop.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        prop.resource_name)

    def should_have_properties(self, properties_list):
        if type(properties_list) is not list:
            properties_list = [properties_list]

        for prop in self.properties:
            property_names = prop.property_value.keys()
            for required_property_name in properties_list:
                if required_property_name not in property_names:
                    self.validator.preprocessor.add_failure("[{0}.{1}.{2}] should have property: '{3}'".format(prop.resource_type,
                                                                                                               prop.resource_name,
                                                                                                               prop.property_name,
                                                                                                               required_property_name),
                                                            prop.moduleName,
                                                            prop.fileName,
                                                            self.validator.severity,
                                                            self.validator.isRuleOverridden,
                                                            self.validator.overrides,
                                                            prop.resource_name)

    def should_not_have_properties(self, properties_list):
        if type(properties_list) is not list:
            properties_list = [properties_list]

        for prop in self.properties:
            property_names = prop.property_value.keys()
            for excluded_property_name in properties_list:
                if excluded_property_name in property_names:
                    self.validator.preprocessor.add_failure("[{0}.{1}.{2}] should not have property: '{3}'".format(prop.resource_type,
                                                                                                                   prop.resource_name,
                                                                                                                   prop.property_name,
                                                                                                                   excluded_property_name),
                                                            prop.moduleName,
                                                            prop.fileName,
                                                            self.validator.severity,
                                                            self.validator.isRuleOverridden,
                                                            self.validator.overrides,
                                                            prop.resource_name)

    def find_property(self, regex):
        lst = TerraformPropertyList(self.validator)
        for prop in self.properties:
            for nested_property in prop.property_value:
                if self.validator.matches_regex_pattern(nested_property, regex):
                    lst.properties.append(TerraformProperty(prop.resource_type,
                                           "{0}.{1}".format(prop.resource_name, prop.property_name),
                                           nested_property,
                                           prop.property_value[nested_property],
                                           prop.moduleName,
                                           prop.fileName))
        return lst

    def should_match_regex(self, regex):
        for prop in self.properties:
            if not self.validator.matches_regex_pattern(prop.property_value, regex):
                self.validator.preprocessor.add_failure("[{0}.{1}] should match regex '{2}'".format(prop.resource_type, "{0}.{1}".format(prop.resource_name, prop.property_name), regex),
                                                        prop.moduleName,
                                                        prop.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        prop.resource_name)

    def should_contain_valid_json(self):
        for prop in self.properties:
            try:
                json.loads(prop.property_value)
            except:
                self.validator.preprocessor.add_failure("[{0}.{1}.{2}] is not valid json".format(prop.resource_type, prop.resource_name, prop.property_name),
                                                        prop.moduleName,
                                                        prop.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        prop.resource_name)

    def bool2str(self, b):
        if str(b).lower() in ["true"]:
            return "True"
        if str(b).lower() in ["false"]:
            return "False"
        return b

    def int2str(self, property_value):
        if type(property_value) is int:
            property_value = str(property_value)
        return property_value


class TerraformProperty:

    def __init__(self, resource_type, resource_name, property_name, property_value, moduleName, fileName):
        self.resource_type = resource_type
        self.resource_name = resource_name
        self.property_name = property_name
        self.property_value = property_value
        self.moduleName = moduleName
        self.fileName = fileName


class TerraformResource:

    def __init__(self, typ, name, config, fileName, moduleName):
        self.type = typ
        self.name = name
        self.config = config
        self.fileName = fileName
        self.moduleName = moduleName


class TerraformResourceList:

    def __init__(self, validator, requestedResourceType, resourceTypes, resources):
        self.validator = validator
        self.resource_list = []
        self.requestedResourceType = requestedResourceType

        resourcesByType = {}
        for resourceName in resources:
            resource = resources[resourceName]
            resourceType = resource.type
            resourcesByType[resourceType] = resourcesByType.get(resourceType, {})
            resourcesByType[resourceType][resourceName] = resource.config

        if type(requestedResourceType) is str:
            resourceTypes = []
            for resourceType in resourcesByType:
                if validator.matches_regex_pattern(resourceType, requestedResourceType):
                    resourceTypes.append(resourceType)
        elif requestedResourceType is not None:
            resourceTypes = requestedResourceType

        for resourceType in resourceTypes:
            if resourceType in resourcesByType.keys():
                for resourceName in resourcesByType[resourceType]:
                    self.resource_list.append(
                        TerraformResource(resourceType, resourceName, resourcesByType[resourceType][resourceName], resources[resourceName].fileName, resources[resourceName].moduleName))

        self.resource_types = resourceTypes

    def property(self, property_name):
        lst = TerraformPropertyList(self.validator)
        if len(self.resource_list) > 0:
            for resource in self.resource_list:
                if property_name in resource.config.keys():
                    lst.properties.append(TerraformProperty(resource.type, resource.name, property_name, resource.config[property_name], resource.moduleName, resource.fileName))
                elif self.validator.raise_error_if_property_missing:
                    self.validator.preprocessor.add_failure("[{0}.{1}] should have property: '{2}'".format(resource.type, resource.name, property_name),
                                                            resource.moduleName,
                                                            resource.fileName,
                                                            self.validator.severity,
                                                            self.validator.isRuleOverridden,
                                                            self.validator.overrides,
                                                            resource.name)

        return lst

    def find_property(self, regex):
        lst = TerraformPropertyList(self.validator)
        if len(self.resource_list) > 0:
            for resource in self.resource_list:
                for prop in resource.config:
                    if self.validator.matches_regex_pattern(prop, regex):
                        lst.properties.append(TerraformProperty(resource.type,
                                                                resource.name,
                                                                prop,
                                                                resource.config[prop],
                                                                resource.moduleName,
                                                                resource.fileName))
        return lst

    def with_property(self, property_name, regex):
        lst = TerraformResourceList(self.validator, None, self.resource_types, {})

        if len(self.resource_list) > 0:
            for resource in self.resource_list:
                for prop in resource.config:
                    if prop == property_name:
                        tf_property = TerraformProperty(resource.type, resource.name, property_name, resource.config[property_name], resource.moduleName, resource.fileName)
                        if self.validator.matches_regex_pattern(tf_property.property_value, regex):
                            lst.resource_list.append(resource)

        return lst

    def should_not_exist(self):
        for terraformResource in self.resource_list:
            if terraformResource.type == self.requestedResourceType:
                self.validator.preprocessor.add_failure("[{0}] should not exist. Found in resource named {1}".format(self.requestedResourceType, terraformResource.name),
                                                        terraformResource.moduleName,
                                                        terraformResource.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        terraformResource.name)

    def should_have_properties(self, properties_list):
        if type(properties_list) is not list:
            properties_list = [properties_list]

        if len(self.resource_list) > 0:
            for resource in self.resource_list:
                property_names = resource.config.keys()
                for required_property_name in properties_list:
                    if required_property_name not in property_names:
                        self.validator.preprocessor.add_failure("[{0}.{1}] should have property: '{2}'".format(resource.type, resource.name, required_property_name),
                                                                resource.moduleName,
                                                                resource.fileName,
                                                                self.validator.severity,
                                                                self.validator.isRuleOverridden,
                                                                self.validator.overrides,
                                                                resource.name)

    def should_not_have_properties(self, properties_list):
        if type(properties_list) is not list:
            properties_list = [properties_list]

        if len(self.resource_list) > 0:
            for resource in self.resource_list:
                property_names = resource.config.keys()
                for excluded_property_name in properties_list:
                    if excluded_property_name in property_names:
                        self.validator.preprocessor.add_failure("[{0}.{1}] should not have property: '{2}'".format(resource.type, resource.name, excluded_property_name),
                                                                resource.moduleName,
                                                                resource.fileName,
                                                                self.validator.severity,
                                                                self.validator.isRuleOverridden,
                                                                self.validator.overrides,
                                                                resource.name)

    def name_should_match_regex(self, regex):
        for resource in self.resource_list:
            if not self.validator.matches_regex_pattern(resource.name, regex):
                self.validator.preprocessor.add_failure("[{0}.{1}] name should match regex: '{2}'".format(resource.type, resource.name, regex),
                                                        resource.moduleName,
                                                        resource.fileName,
                                                        self.validator.severity,
                                                        self.validator.isRuleOverridden,
                                                        self.validator.overrides,
                                                        resource.name)


class Validator:

    # default severity is high
    severity = "high"
    preprocessor = None

    def __init__(self):
        self.raise_error_if_property_missing = False

    def resources(self, typ):
        resources = self.terraform.get('resource', {})

        return TerraformResourceList(self, typ, None, resources)

    def error_if_property_missing(self):
        self.raise_error_if_property_missing = True

    # generator that loops through all files to be scanned (stored internally in fileName; returns self (Validator) but sets self.fileName and self.terraform
    def get_terraform_files(self, isRuleOverridden):
        self.isRuleOverridden = isRuleOverridden
        for self.fileName, self.terraform in self.preprocessor.modulesDict.items():
            yield self

    def matches_regex_pattern(self, variable, regex):
        return not (self.get_regex_matches(regex, variable) is None)

    def get_regex_matches(self, regex, variable):
        if regex[-1:] != "$":
            regex = regex + "$"

        if regex[0] != "^":
            regex = "^" + regex

        variable = str(variable)
        if '\n' in variable:
            return re.match(regex, variable, re.DOTALL)
        return re.match(regex, variable)

    # this is only used by unit_test.py
    def setTerraform(self, terraform):
        self.terraform = terraform
        self.fileName = "none.tf"


class PreProcessor:

    TF = ".tf"
    UTF8 = "utf8"
    IS_MODULE = "__isModule__"
    PARENT = "__parent__"
    MODULE_NAME = "__ModuleName__"
    FILE_NAME = "__fileName__"
    LOCALS = "locals"
    VARIABLE = "variable"
    OUTPUT = "output"
    RESOURCE = "resource"
    VALUE = "value"
    MODULE = "module"
    SOURCE = "source"
    DEFAULT = "default"
    REGEX_COLON_BRACKET = re.compile('.*:\s*\[.*', re.DOTALL)           # any characters : whitespace [ any characters

    def __init__(self, jsonOutput):
        self.jsonOutput = jsonOutput
        self.variablesFromCommandLine = {}
        self.hclDict = {}
        self.modulesDict = {}
        self.fileNames = {}
        self.passNumber = 1
        self.dummyIndex = 0
        # on 1st pass replace var. with var$.
        # on 2nd pass replace var. and var$. with var!.
        self.braces = ["${", "@{", "!{"]
        self.vars = ["var.", "var@.", "var!."]
        self.locals = ["local.", "local@.", "local!."]
        self.modules = ["module.", "module@.", "module!."]
        self.terraform_workspaces = ["terraform.workspace", "terraform@.workspace", "terraform!.workspace"]
        self.datas = ["data.", "data@.", "data!."]
        self.variableFind = [self.braces[0], self.vars[0], self.locals[0], self.modules[0], self.terraform_workspaces[0], self.datas[0]]
        self.variableErrorReplacement = [self.braces[1], self.vars[1], self.locals[1], self.modules[1], self.terraform_workspaces[1], self.datas[1]]
        self.variableErrorReplacementPass2 = [self.braces[2], self.vars[2], self.locals[2], self.modules[2], self.terraform_workspaces[2], self.datas[2]]
        self.replacements = [self.variableFind, self.variableErrorReplacement, self.variableErrorReplacementPass2]
        self.replaceableVariablePrefixes = [self.braces[0], self.vars[0], self.locals[0], self.modules[0], self.terraform_workspaces[0],
                                            self.braces[1], self.vars[1], self.locals[1], self.modules[1], self.terraform_workspaces[1],
                                            self.braces[2], self.vars[2], self.locals[2], self.modules[2], self.terraform_workspaces[2]]
        
    def process(self, path, variablesJsonFilename=None):
        inputVars = {}
        if variablesJsonFilename is not None:
            for fileName in variablesJsonFilename:
                with open(fileName, "r", encoding="utf-8") as fp:
                    try:
                        variables_string = fp.read()
                        inputVarsDict = hcl.loads(variables_string)
                        inputVars = {**inputVars, **inputVarsDict}
                    except:
                        self.add_error_force(traceback.format_exc(), "---", fileName, "high")

        # prefix any input variable not containing '.' with 'var.'
        for var in inputVars:
            if "." not in var:
                newVar = "var." + var
                self.variablesFromCommandLine[newVar] = inputVars[var]
            else:
                self.variablesFromCommandLine[var] = inputVars[var]
                
        self.root = None
        self.readDir(path, self.hclDict)

        # all terraform files are now loaded into hclDict (indexed by subdirectory/fileName/terraform structure)
        # process hclDict and load every module into modulesDict
        self.getAllModules(self.hclDict, False)
        # make second pass so variables depending on a previous definition should now be defined
        logging.warning("------------------>>>starting pass 2...")
        self.passNumber = 2
        self.variableFind += self.variableErrorReplacement
        self.variableErrorReplacement = self.variableErrorReplacementPass2 + self.variableErrorReplacementPass2
        self.getAllModules(self.hclDict, False)

    def readDir(self, path, d):
        for directory, subdirectories, files in os.walk(path):
            if self.root is None:
                self.root = directory
                i = self.root.rfind(os.path.sep)
                if i != -1:
                    self.root = self.root[i+1:]
                # define root and mark this dictionary as a module since all directories in terraform are modules by default
                self.hclDict[self.root] = {}
                self.hclDict[self.root][self.IS_MODULE] = True
                self.hclDict[self.root][self.PARENT] = None
                self.hclDict[self.root][self.MODULE_NAME] = self.root
                d = self.hclDict[self.root]

            for file in files:
                if file[-3:].lower() == self.TF:
                    # terraform file (ends with .tf)
                    fileName = os.path.join(directory, file)
                    relativeFileName = fileName[len(path):]
                    with open(fileName, 'r', encoding='utf8') as fp:
                        try:
                            terraform_string = fp.read()
                            if len(terraform_string.strip()) > 0:
                                self.loadFileByDir(fileName, relativeFileName, d, d, terraform_string)
                                self.fileNames[fileName] = fileName
                        except:
                            self.add_error_force(traceback.format_exc(), "---", fileName, "high")

    # load file by directory, marking each directory as a module and setting parent directories
    def loadFileByDir(self, fileName, path, hclSubDirDict, parentDir, terraform_string):
        i = path.find("\\")
        if i == -1:
            # \ not found, try /
            i = path.find("/")
            if i == -1:
                # end of subdirectories; path is a terraform filename; load terraform file into dictionary
                hclSubDirDict[path] = hcl.loads(terraform_string)
                hclSubDirDict[path][self.FILE_NAME] = fileName
                # remove file name from end of path
                t = self.getPreviousLevel(fileName, os.path.sep)
                self.findModuleSources(hclSubDirDict[path], parentDir, t[0])
                return
        if i == 0:
            # found in first character, recursively try again skipping first character
            self.loadFileByDir(fileName, path[1:], hclSubDirDict, parentDir, terraform_string)
        else:
            # get subdirectory
            subdir = path[:i]
            if hclSubDirDict.get(subdir) is None:
                # subdirectory not defined in our dictionary yet so define it
                hclSubDirDict[subdir] = {}
                hclSubDir = hclSubDirDict[subdir]
                # mark this dictionary as a module since all directories in terraform are modules by default
                hclSubDir[self.IS_MODULE] = True
                hclSubDir[self.PARENT] = parentDir
                hclSubDir[self.MODULE_NAME] = subdir
            else:
                hclSubDir = hclSubDirDict[subdir]
            # recursively process next subdirectory
            i += 1
            self.loadFileByDir(fileName, path[i:], hclSubDirDict[subdir], hclSubDir, terraform_string)

    def findModuleSources(self, d, parentDir, currentFileName):
        for key in d:
            # only process module key
            if key == self.MODULE:
                modules = d[key]
                # process all modules
                for moduleName in modules:
                    module = modules[moduleName]
                    # find source parameter
                    for parameter in module:
                        if parameter == self.SOURCE:
                            sourcePath = module[parameter]
                            if not sourcePath.startswith("git::"):
                                self.createMissingFromSourcePath(sourcePath, parentDir, currentFileName)

    def createMissingFromSourcePath(self, sourcePath, d, currentFileName):
        # source is local
        while sourcePath != "":
            t = self.getNextLevel(sourcePath, "/")
            currentModule = t[0]
            sourcePath = t[1]
            if currentModule == "..":
                # move up a level
                t = self.getPreviousLevel(currentFileName, os.path.sep)
                currentFileName = t[0]
                if d.get(self.PARENT) is None:
                    # add parent
                    self.dummyIndex += 1
                    parent = {}
                    parent[self.PARENT] = None
                    parent[self.IS_MODULE] = True
                    parent[self.MODULE_NAME] = "dummy" + str(self.dummyIndex)
                    parent[d[self.MODULE_NAME]] = d
                    d[self.PARENT] = parent
                    self.hclDict = parent
                d = d[self.PARENT]
            elif currentModule == ".":
                # current directory; do nothing
                pass
            else:
                # move down to currentModule level
                currentFileName += os.path.sep + currentModule
                md = d.get(currentModule, False)
                if md is False:
                    # create new level
                    d[currentModule] = {}
                    md = d[currentModule]
                    md[self.PARENT] = d
                    md[self.MODULE_NAME] = currentModule
                    md[self.IS_MODULE] = True

                d = md
        # read directory
        self.readDir(currentFileName, d)
        return d

    # get all modules for given dictionary d; pass isModule as True if in a module block
    def getAllModules(self, d, isModule):
        for key in d:
            # ignore parent key
            if key != self.PARENT:
                value = d[key]
                if type(value) is dict:
                    if isModule or self.isModule(value):
                        moduleName = key
                        # load module, resolve variables in it and add it to modules dictionary
                        moduleDict = self.getModule(moduleName)
                        if self.isModule(value):
                            moduleDict[self.PARENT] = d[moduleName][self.PARENT]
                    # recursively get all modules in the nested dictionary in value
                    self.getAllModules(value, key == "module")

    # get given moduleName from modulesDict; find & load it if not there yet
    def getModule(self, moduleName, errorIfNotFound=True, dictToCopyFrom=None, tfDict=None):
        moduleDict = self.modulesDict.get(moduleName)
        if moduleDict is None or (moduleDict[self.VARIABLE] == {} and moduleDict[self.LOCALS] == {} and moduleDict[self.OUTPUT] == {}):
            # not there yet, find it and load it
            moduleDict = self.findModule(moduleName, self.hclDict, dictToCopyFrom, tfDict)
            if moduleDict is None:
                # couldn't find it, log it and create a dummy entry
                if errorIfNotFound:
                    self.logMsg("error", "Couldn't find module " + moduleName)
                moduleDict = self.createModuleEntry(moduleName)
        elif self.passNumber > 1:
            # module found on second pass, re-resolve variables
            self.findModule(moduleName, self.hclDict, dictToCopyFrom, tfDict)

        return moduleDict

    # find given moduleName in given dictionary d; load module attributes and resolve variables in module; last two parameters are all or nothing
    def findModule(self, moduleName, d, dictToCopyFrom=None, tfDict=None):
        # use dictToCopyFrom if provided
        if dictToCopyFrom is not None:
            sourcePath = self.getSourcePath(dictToCopyFrom)
            if sourcePath is not None and not sourcePath.startswith("git::"):
                dd = self.getModuleDictFromSourcePath(sourcePath, tfDict)
                if dd:
                    moduleDict = self.modulesDict.get(dd[self.MODULE_NAME])
                    if moduleDict is None:
                        moduleDict = self.createModuleEntry(dd[self.MODULE_NAME])
                        moduleDict[self.IS_MODULE] = True
                    self.loadModule(moduleName, dd, dictToCopyFrom)
                    if moduleName != dd[self.MODULE_NAME]:
                        self.loadModule(dd[self.MODULE_NAME], dd, dictToCopyFrom)
                    # source module found, replace variables and return it
                    m = self.loadModule(dd[self.MODULE_NAME], dd, dictToCopyFrom)
                    return m
            return self.loadModule(moduleName, {}, dictToCopyFrom)

        for key in d:
            # ignore parent key
            if key != self.PARENT:
                value = d[key]
                if key == moduleName:
                    # module found, replace variables and return it
                    return self.loadModule(moduleName, value, dictToCopyFrom)
                else:
                    if type(value) is dict:
                        if self.isModule(value):
                            # recursively find the module
                            m = self.findModule(moduleName, value)
                            if m is not None:
                                return m
        # not found
        return None

    def loadModule(self, moduleName, d, dictToCopyFrom):
        self.logMsgAlways("warning", ">>>loading module " + moduleName)

        moduleDict = self.modulesDict.get(moduleName)
        if moduleDict is None:
            # create empty module entry
            moduleDict = self.createModuleEntry(moduleName)
            moduleDict[self.IS_MODULE] = self.hasTerraform(d)

        if dictToCopyFrom is not None:
            mdv = moduleDict[self.VARIABLE]
            # add/replace the passed in variables to the module's variables
            for attr in dictToCopyFrom:
                if attr != self.SOURCE:
                    # only replace on pass #1 if resolved
                    if self.passNumber == 2 or self.isResolved(dictToCopyFrom[attr]):
                        mdv[attr] = dictToCopyFrom[attr]

        # load all attributes for this module
        self.loadModuleAttributes(moduleName, d, moduleDict, None)
        # resolve variables for this module
        self.resolveVariablesInModule(moduleName, moduleDict)
        return moduleDict

    def isResolved(self, var):
        if type(var) is str:
            return self.isStrResolved(var)
        elif type(var) is dict:
            for key in var:
                if not self.isResolved(var[key]):
                    return False
        elif type(var) is list:
            for value in var:
                if not self.isResolved(value):
                    return False
        else:
            return False
        return True

    def isStrResolved(self, var):
        for varErrorReplacement in self.variableErrorReplacement:
            if varErrorReplacement in var:
                return False
            
        return True

    def createModuleEntry(self, moduleName):
        self.modulesDict[moduleName] = {}
        moduleDict = self.modulesDict[moduleName]
        moduleDict[self.VARIABLE] = {}
        moduleDict[self.LOCALS] = {}
        moduleDict[self.OUTPUT] = {}
        moduleDict[self.RESOURCE] = {}
        moduleDict[self.IS_MODULE] = False
        return moduleDict

    def loadModuleAttributes(self, moduleName, d, moduleDict, tfDict):
        if self.isModule(d):
            if tfDict is None:
                tfDict = d
            else:
                # skip nested modules
                return

        for key in sorted(d):
            # ignore parent key
            if key != self.PARENT:
                value = d[key]
                if key == self.LOCALS:
                    # get values for all local variables
                    for local in value:
                        # only replace on first pass or not already fully resolved
                        if self.passNumber == 1 or self.containsVariable(moduleDict[self.LOCALS][local]):
                            moduleDict[self.LOCALS][local] = value[local]
                elif key == self.OUTPUT:
                    for output in value:
                        # only replace on first pass or not already fully resolved
                        if self.passNumber == 1 or self.containsVariable(moduleDict[self.OUTPUT][output]):
                            moduleDict[self.OUTPUT][output] = value[output][self.VALUE]
                elif key == self.RESOURCE:
                    for resourceType in value:
                        resourceNames = value[resourceType]
                        for resourceName in resourceNames:
                            config = resourceNames[resourceName]
                            res = moduleDict[self.RESOURCE].get(resourceName, None)
                            # only replace on first pass or not already fully resolved
                            if self.passNumber == 1 or (res != None and self.containsVariable(res.config)):
                                moduleDict[self.RESOURCE][resourceName] = TerraformResource(resourceType, resourceName, config, d[self.FILE_NAME], moduleName)
                elif key == self.VARIABLE:
                    '''
                    value could be a string as in below case
                        condition {
                            test = "ArnEquals"
                            variable = "aws:SourceArn"
                            values = ["${var.services_entry_arn}"]
                        }
                    '''
                    if type(value) is dict:
                        # initialize any default values for variables
                        for variable in value:
                            if value[variable].get(self.DEFAULT) is not None:
                                moduleDict[self.VARIABLE][variable] = value[variable][self.DEFAULT]
                elif key == self.MODULE:
                    # loop through all modules
                    for mn in value:
                        # resolve parameter variables first
                        for parameter in value[mn]:
                            if parameter != self.SOURCE:
                                replacementValue = self.resolveVariableByType(value[mn][parameter], moduleName)
                                if replacementValue != value[mn][parameter]:
                                    self.logMsgAlways("warning", "replaced module " + mn + " parameter " + parameter + " value " + str(value[mn][parameter]) + " with " + str(replacementValue))
                                    value[mn][parameter] = replacementValue
                        # get defined module; load it if not already there
                        md = self.getModule(mn, False, value[mn], tfDict)
                        # copy all outputs from source module (md) to containing module variable
                        for output in md[self.OUTPUT]:
                            # only copy if not already there
                            if moduleDict[self.VARIABLE].get(output) is None:
                                moduleDict[self.VARIABLE][output] = md[self.OUTPUT][output]
                else:
                    if type(value) is dict:
                        # don't load any other nested modules
                        if not self.isModule(value):
                            self.loadModuleAttributes(moduleName, value, moduleDict, tfDict)

    def getSourcePath(self, parameterDict):
        for parameter in parameterDict:
            if parameter == self.SOURCE:
                return parameterDict[parameter]
        return None

    def getModuleDictFromSourcePath(self, sourcePath, d):
        # source is local
        while sourcePath != "":
            t = self.getNextLevel(sourcePath, "/")
            currentModule = t[0]
            sourcePath = t[1]
            if currentModule == "..":
                # move up a level
                d = d[self.PARENT]
            elif currentModule == ".":
                # current directory; do nothing
                pass
            else:
                # move down to currentModule level
                d = d.get(currentModule, False)
                if d is False:
                    return False
        return d

    # resolve variables (anything surrounded by ${}) in given moduleDict
    def resolveVariablesInModule(self, moduleName, moduleDict):
        self.shouldLogErrors = False
        # resolve variables
        for key in moduleDict[self.VARIABLE]:
            value = moduleDict[self.VARIABLE][key]
            replacementValue = self.resolveVariableByType(value, moduleName)
            moduleDict[self.VARIABLE][key] = replacementValue
            if replacementValue != value:
                self.logMsgAlways("warning", "replaced variable " + key + " value " + str(value) + " with " + str(replacementValue))
        # resolve locals
        for key in moduleDict[self.LOCALS]:
            value = moduleDict[self.LOCALS][key]
            replacementValue = self.resolveVariableByType(value, moduleName)
            moduleDict[self.LOCALS][key] = replacementValue
            if replacementValue != value:
                self.logMsgAlways("warning", "replaced local variable " + key + " value " + str(value) + " with " + str(replacementValue))
        # resolve outputs
        for key in moduleDict[self.OUTPUT]:
            value = moduleDict[self.OUTPUT][key]
            replacementValue = self.resolveVariableByType(value, moduleName)
            moduleDict[self.OUTPUT][key] = replacementValue
            if replacementValue != value:
                self.logMsgAlways("warning", "replaced output variable " + key + " value " + str(value) + " with " + str(replacementValue))
        # resolve resources
        self.shouldLogErrors = True
        for key in moduleDict[self.RESOURCE]:
            value = moduleDict[self.RESOURCE][key].config
            replacementValue = self.resolveVariableByType(value, moduleName)
            moduleDict[self.RESOURCE][key].config = replacementValue
            if replacementValue != value:
                self.logMsgAlways("warning", "replaced resource variable " + key + " value " + str(value) + " with " + str(replacementValue))

    def resolveVariableByType(self, value, moduleName):
        if type(value) is str:
            return self.resolveVariableLine(value, moduleName)
        elif type(value) is dict:
            return self.resolveDictVariable(value, moduleName)
        elif type(value) is list:
            return self.resolveListVariable(value, moduleName)
        elif type(value) is tuple:
            return self.resolveTupleVariable(value, moduleName)
        else:
            return value

    def resolveDictVariable(self, value, moduleName):
        returnValue = {}
        for key in value:
            returnValue[key] = self.resolveVariableByType(value[key], moduleName)
        return returnValue

    def resolveListVariable(self, value, moduleName):
        if len(value) == 0:
            return value;
        index = 0
        for v in value:
            value[index] = self.resolveVariableByType(v, moduleName)
            index += 1
        if value[0] in ("join", "merge", "concat", "coalesce", "element", "coalescelist"):
            # supported function
            return self.handleFunction(value)
        return value

    def resolveTupleVariable(self, value, moduleName):
        returnValue = tuple(self.resolveVariableByType(v, moduleName) for v in value)
        if len(returnValue) == 3:
            floatValue0 = self.getFloatValue(returnValue[0])
            floatValue2 = self.getFloatValue(returnValue[2])
            if type(floatValue0) == float and type(returnValue[1]) == str and type(floatValue2) == float:
                if returnValue[1] == "+":
                    return floatValue0 + floatValue2
                elif returnValue[1] == "-":
                    return floatValue0 - floatValue2
                elif returnValue[1] == "*":
                    return floatValue0 * floatValue2
                elif returnValue[1] == "/":
                    return floatValue0 / floatValue2
        return returnValue
            
    def getFloatValue(self, value):
        try:
            return float(value)
        except:
            return value;

    def handleFunction(self, value):
        # check if all variables have been resolved
        if self.containsVariable(value, True):
            # not fully resolved yet; return what we have so far
            return value
        return self.processFunction(value)
        
    def processFunction(self, value):
        it = iter(value)
        function = next(it, None)
        if function == "join":
            delimiter = next(it, None)
            t = next(it, None)
            return delimiter.join(v for v in t)
        elif function == "merge":
            d = {}
            for v in it:
                if type(v) is dict:
                    for key in v:
                        d[key] = v[key]
                else:
                    d[v] = v
            return d;
        elif function == "concat":
            d = []
            for v in it:
                if type(v) is list:
                    for entry in v:
                        d.append(entry)
                else:
                    d.append(v)
            return d;
        elif function == "element":
            lst = next(it, None)
            index = next(it, None)
            if '*' in lst:
                return value
            return lst[index]
        else:
            # coalesce/coalescelist
            d = {}
            if value[len(value)-1] == "...":
                # there is a single list that needs to be processed
                for v in value[1]:
                    if v:
                        if type(v) is tuple:
                            return self.processFunction(v)
                        else:
                            return v
            else:
                for v in it:
                    if v:
                        if type(v) is tuple:
                            return self.processFunction(v)
                        else:
                            return v
            # no non-empty entries; undefined what to do so return None
            return None;
        

    # returns True if given dictionary d contains a key of __isModule__
    def isModule(self, d):
        for key in d:
            if key == self.IS_MODULE:
                return d[key]
        return False

    # returns True if given dictionary d contains at least one terraform file
    def hasTerraform(self, d):
        for key in d:
            if key.lower().endswith(self.TF):
                return True
        return False

    # resolve entire variable
    def resolveVariableLine(self, value, moduleName):
        if not self.containsVariable(value):
            return value
        # a variable needs to be replaced
        t = self.findVariable(value, True)
        var = t[0]
        b = t[1]
        e = t[2]
        if var.startswith("["):
            var = var[1:len(var)-1]
        var = var.strip()
        rv = self.resolveVariable(var, moduleName)
        if b == 0 and (e == len(value) or e == -1):
            # full replacement; don't merge since a string may not have been returned
            newValue = rv[0]
        else:
            newValue = value[:b] + str(rv[0]) + value[e:]
        # recursively resolve the variables since there may be more than one variable in this value
        return self.resolveVariableByType(newValue, moduleName)

    # resolve innermost variable
    def resolveVariable(self, value, moduleName, dictToCopyFrom=None, tfDict=None):
        # find variable (possibly in brackets)
        isOldTFvarStyle=False
        v, b, e, insideBrackets, foundDelineator, foundDelineatorErrRepl = self.findVariable(value, False)
        if len(v) > 1 and v[1] == "{":
            isOldTFvarStyle = True
        if not insideBrackets and isOldTFvarStyle:
            # inside ${}; remove them
            var = value[2:e-1]
        else:
            var = v
        # update moduleName in case we switch modules and need to recurse more
        replacementValue, moduleName, isHandledType = self.getReplacementValue(var, moduleName, isOldTFvarStyle, dictToCopyFrom, tfDict)
        if replacementValue == var:
            # couldn't find a replacement; change to our notation to mark it
            if isHandledType:
                self.logMsg("error", "Couldn't find a replacement for: " + self.getOrigVar(var) + " in " + moduleName)
            else:
                self.logMsg("debug", "Couldn't find a replacement for: " + self.getOrigVar(var) + " in " + moduleName)
            if not isOldTFvarStyle:
                # strip off replaceable variable
                var = var[len(foundDelineator):]
            replacementValue = value[:b] + foundDelineatorErrRepl + var
            if insideBrackets:
                replacementValue += "]"
            if not isOldTFvarStyle and len(value) > 1 and value[1] == "{":
                replacementValue += "}"
            if isOldTFvarStyle and e > 0:
                # remove closing brace
                replacementValue += value[e-1:]
            return (replacementValue, not insideBrackets)

        if type(replacementValue) is str:
            if insideBrackets:
                self.logMsgAlways("info", "  replacing [" + var + "] with " + replacementValue)
                # resolve the variable again since the replacement may also contain variables
                return (value[:b] + self.resolveVariableLine(replacementValue, moduleName) + value[e:], not insideBrackets)
            else:
                if v == replacementValue:
                    # this prevents a loop
                    replacementValue = replacementValue.replace(foundDelineator, foundDelineatorErrRepl, 1)
                    self.logMsg("debug", "Couldn't find a replacement for: " + self.getOrigVar(var) + " (would have looped) in " + moduleName)
                else:
                    self.logMsgAlways("info", "  replacing ${" + var + "} with " + replacementValue)
                    # resolve the variable again since the replacement may also contain variables
                    return (self.resolveVariableLine(replacementValue, moduleName), not insideBrackets)
        else:
            if isOldTFvarStyle:
                self.logMsgAlways("info", "  replacing " + foundDelineator + var + "} with " + str(replacementValue))
            else:
                self.logMsgAlways("info", "  replacing " + var + " with " + str(replacementValue))

        return (replacementValue, not insideBrackets)

    def getOrigVar(self, var):
        if var.startswith(self.vars[1]) or var.startswith(self.vars[2]):
            return self.vars[0] + var[len(self.vars[1]):]
        elif var.startswith(self.locals[1]) or var.startswith(self.locals[2]):
            return self.locals[0] + var[len(self.locals[1]):]
        elif var.startswith(self.modules[1]) or var.startswith(self.modules[2]):
            return self.modules[0] + var[len(self.modules[1]):]
        else:
            return var

    # check if given value contains a variable anywhere
    def containsVariable(self, value, isAnyVar=False):
        if type(value) is str:
            t = self.findAnyVariableDelineatorsForVars(value, False, isAnyVar)
            if t[0] == -1:
                return False
            return True
        elif type(value) is dict:
            return self.containsVariableDict(value, isAnyVar)
        elif type(value) is list:
            return self.containsVariableList(value, isAnyVar)
        else:
            # not a variable
            return False

    def containsVariableDict(self, value, isAnyVar):
        for key in value:
            if self.containsVariable(value[key], isAnyVar):
                return True
        # no variables found
        return False

    def containsVariableList(self, value, isAnyVar):
        for v in value:
            if self.containsVariable(v, isAnyVar):
                return True
        # no variables found
        return False

    # find deepest nested variable in given value
    def findVariable(self, value, isNested, previouslyFoundVar=None):
        # pass 1: if unreplaceable, change $ to @
        # pass 2: if unreplaceable, change both $ & @ to !
        if type(value) is str:
            isVar = False
            val = value
            if previouslyFoundVar:
                insideBrackets = previouslyFoundVar[3]
            else:
                insideBrackets = False
            if isNested and type(previouslyFoundVar) is str and "{" in previouslyFoundVar[0]:
                # if this is a nested call and the outer call found ${, only look for the brace now
                braceOnly = True
            else:
                braceOnly = False
 
            b, e, foundDelineator, foundDelineatorErrRepl = self.findVariableDelineatorsForVars(val, braceOnly, self.variableFind, self.variableErrorReplacement)
            if b == -1:
                return None                
            if b > 0:
                partial = value[:b]
                if partial[len(partial)-1] == "[":
                    # open bracket found before the variable
                    insideBrackets = True
                    partial = value[b:e]
                    if partial[len(partial)-1] == "]":
                        # close bracket found after the variable, remove from variable end
                        e -= 1
                
            isVar = True
            if "{" in foundDelineator and e == -1:
                # problem
                self.add_error("Matching close brace not found: " + value, "---", "---", "high")
                return None
            
            foundVar = (value[b:e], b, e, insideBrackets, foundDelineator, foundDelineatorErrRepl)
            
            newSearchValue = foundVar[0]
            bOffset = 0
            if isVar:
                # remove delineator(s)
                if newSearchValue.endswith("}"):
                    newSearchValue = newSearchValue[2:len(newSearchValue)-1]
                    bOffset = 2
                else:
                    newSearchValue = newSearchValue[len(foundDelineator):]
                    bOffset = len(foundDelineator)
            else:
                newSearchValue = newSearchValue[1:len(newSearchValue)-1]
                if not insideBrackets:
                    # adjust beginning & ending since nested inside previouslyFoundVar
                    fv_length = len(previouslyFoundVar[4])
                    b += fv_length
                    e += fv_length
                    bOffset = 0
                else:
                    bOffset = 1
                foundVar = (newSearchValue, b, e, insideBrackets, foundDelineator, foundDelineatorErrRepl)
            if newSearchValue not in self.terraform_workspaces[0]:
                # recursively find variable
                fv = self.findVariable(newSearchValue, True, foundVar)
                if fv is None:
                    # no variable found
                    return foundVar
                if foundVar[0].endswith("}") and fv[1] == 0 and fv[2] == len(newSearchValue):
                    # return originally found variable which is the old style
                    return foundVar
                if insideBrackets:
                    # use beginning & ending from previous
                    fv = (fv[0], b, e, fv[3], fv[4], fv[5])
                else:
                    fv = (fv[0], fv[1]+b+bOffset, fv[2]+b+bOffset, fv[3], fv[4], fv[5])
                return fv
            else:
                return foundVar
        return previouslyFoundVar

    def findAnyVariableDelineatorsForVars(self, value, braceOnly, isAnyVar):
        if isAnyVar:
            for variableFind in self.replaceableVariablePrefixes:
                t = self.findVariableDelineatorsForVars(value, braceOnly, [variableFind], [variableFind])
                if t[0] != -1:
                    return t
            return -1, 0, None, None
        return self.findVariableDelineatorsForVars(value, braceOnly, self.variableFind, self.variableErrorReplacement)
        
    def findVariableDelineatorsForVars(self, value, braceOnly, variableFind, variableErrorReplacement):
        if braceOnly and value not in self.terraform_workspaces:
            b, e = self.findVariableDelineators(value, variableFind[0], "}", variableErrorReplacement[0])
            if b > -1:
                return b, e, variableFind[0], variableErrorReplacement[0]
        else:
            prevB = -1
            for varPrefix, varErrorReplacement in zip(variableFind, variableErrorReplacement):
                if varPrefix[1] == "{":
                    closeVar = "}"
                else:
                    closeVar = None
                b, e = self.findVariableDelineators(value, varPrefix, closeVar, varErrorReplacement)
                if b > -1:
                    if closeVar != None or b == 0 or (value[b-1] != "{"):
                        if b > prevB:
                            prevB = b;
                            prevE = e
                            prevVarPrefix = varPrefix
                            prevVarErrorReplacement = varErrorReplacement
            if prevB != -1:
                return prevB, prevE, prevVarPrefix, prevVarErrorReplacement
        return -1, 0, None, None

    def findVariableDelineators(self, value, openVar, closeVar, varErrorReplacement=None):
        '''
        This is valid:  name-prefix = "sf-${module.common.account_name}-${local.pcas_vpc_type}-${local.env}-${module.common.region}"
        This is not:    name-prefix = "sf-module.common.account_name-local.pcas_vpc_type-local.env-module.common.region"
        i.e. if combining variables into one variable, must use old ${} variable style
        '''
        b = value.rfind(openVar)
        if b == -1:
            return -1, 0
        if openVar == "[":
            # check if preceeded by :
            matchObject = self.REGEX_COLON_BRACKET.search(value)
            if matchObject != None:
                return -1, 0
        if closeVar == None:
            if openVar in self.terraform_workspaces:
                return b, b + len(openVar)
            # search for default "closeVars"
            if value.find("[", b) == -1:
                defaultCloseVars = [",", "'", ")", "}", "]"]
            else:
                # don't include ] if [ is in value
                defaultCloseVars = [",", "'", ")", "}"]
            prevE = 99999
            for closeVar in defaultCloseVars:
                e = value.find(closeVar, b)
                if e != -1 and e < prevE:
                    prevE = e
            if prevE != 99999:
                return b, prevE
            # no close value found, use length of value
            return b, len(value)
        v = value[b+1:]
        nested = 0
        for index, char in enumerate(v):
            if char == closeVar:
                if nested == 0:
                    return b, b+index+2
                # just closed a nested variable
                nested -= 1
            if v[index:index+len(openVar)] == openVar or (varErrorReplacement != None and len(v) >= index+len(varErrorReplacement) and v[index:index+len(varErrorReplacement)] == varErrorReplacement):
                # start of a nested variable
                nested += 1
        # error: matching closeVar not found
        return 0, -1

    # find replacement value for given var in given moduleName
    def getReplacementValue(self, var, moduleName, isOldTFvarStyle, dictToCopyFrom=None, tfDict=None):
        replacementValue = None
        if (var.startswith('"') and var.endswith('"')) or (var.startswith("'") and var.endswith("'")):
            return var[1:len(var)-1], moduleName, True
        subscript = None
        b = var.find("[")
        if b != -1:
            v = var[:b]
            e = var.find("]", b)
            if e != -1:
                subscript = var[b+1:e]
                v += var[e+1:]
                if subscript[0] == '"' or subscript[0] == "'":
                    # remove quotes
                    subscript = subscript[1:len(subscript)-1]
        else:
            v = var

        isHandledType = False
        notHandled = ["?", "==", "!=", ">", "<", ">=", "<=", "&&", "||", "!", "+", "-", "*", "/", "%"]
        moduleDict = self.modulesDict[moduleName]
        if isOldTFvarStyle:
            varsTuple = self.vars[0]
            localsTuple = self.locals[0]
            modulesTuple = self.modules[0]
            terraform_workspacesTuple = self.terraform_workspaces[0]
        else:
            if self.passNumber == 1:
                varsTuple = self.vars[0]
                localsTuple = self.locals[0]
                modulesTuple = self.modules[0]
                terraform_workspacesTuple = self.terraform_workspaces[0]
            else:
                varsTuple = (self.vars[0], self.vars[1])
                localsTuple = (self.locals[0], self.locals[1])
                modulesTuple = (self.modules[0], self.modules[1])
                terraform_workspacesTuple = (self.terraform_workspaces[0], self.terraform_workspaces[1])
        if v.startswith(varsTuple):
            # conditional statements, boolean statements, and math are not currently handled
            if not any(x in v for x in notHandled):
                isHandledType = True
            v = v[self.getPrefixLength(v, varsTuple):]
            index = v.find('.')
            if index > -1:
                subscript = v[index+1:]
                v = v[:index]
            replacementValue = moduleDict[self.VARIABLE].get(v)
        elif v.startswith(localsTuple):
            if not any(x in v for x in notHandled):
                isHandledType = True
            v = v[self.getPrefixLength(v, localsTuple):]
            replacementValue = moduleDict[self.LOCALS].get(v)
        elif v.startswith(modulesTuple):
            if not any(x in v for x in notHandled):
                isHandledType = True
            # variable is in a different module
            modulePrefixLength = self.getPrefixLength(v, modulesTuple)
            e = v.find(".", modulePrefixLength)
            if e == -1:
                self.add_error("Error Resolving module variable: " + var + "  expected ending '.' not found", moduleName, "---", "high")
            else:
                moduleName = v[modulePrefixLength:e]
                md = self.getModule(moduleName, True, dictToCopyFrom, tfDict)
                moduleOutputDict = md[self.OUTPUT]
                e += 1
                remainingVar = v[e:]
                while remainingVar != "":
                    t = self.getNextLevel(remainingVar, ".")
                    moduleOutputDict = moduleOutputDict.get(t[0])
                    if moduleOutputDict is None:
                        self.logMsg("error", "Error resolving variable: " + var + "  variable not found in module (no module source available?) in " + moduleName)
                        return var, moduleName, True
                    remainingVar = t[1]
                if type(moduleOutputDict) is dict:
                    replacementValue = moduleOutputDict.get(self.VALUE, moduleOutputDict)
                else:
                    replacementValue = moduleOutputDict
        elif v.startswith(terraform_workspacesTuple):
            isHandledType = True

        if type(replacementValue) is dict and subscript != None:
            replacementValue = replacementValue.get(subscript, subscript)

        if replacementValue is None:
            replacementValue = self.variablesFromCommandLine.get(var, var)

        return replacementValue, moduleName, isHandledType

    def getPrefixLength(self, var, varTuple):
        if type(varTuple) is str:
            return len(varTuple)
        for t in varTuple:
            if var.startswith(t):
                return len(t)

    def getPreviousLevel(self, var, separator):
        b = var.rfind(separator)
        if b == -1:
            b = len(var)
        return (var[:b], var[b+1:])

    def getNextLevel(self, var, separator):
        b = var.find(separator)
        if b == -1:
            b = len(var)
        return (var[:b], var[b+1:])

    # add given failure in given fileName
    def add_failure(self, failure, moduleName, fileName, severity, isRuleOverridden, overrides, resourceName):
        waived = ""
        waiver = self.overridden(isRuleOverridden, overrides, resourceName, severity)
        if waiver is not None:
            if severity == "high":
                waived = "**waived by " + waiver + "**"
            else:
                waived = "**waived**"
        self.jsonOutput["failures"].append( self.getFailureMsg(severity, waived, failure, moduleName, fileName) )

    def overridden(self, isRuleOverridden, overrides, resourceName, severity):
        if isRuleOverridden:
            for override in overrides:
                if override[1] == resourceName:
                    if severity == "high":
                        if len(override) == 3:
                            return override[2]
                        print("***Invalid override: " + ":".join(override))
                        print("high severity rules must include RR or RAR")
                        print("Needs to be in the following format:  rule_name:resource_name:RR-xxx or rule_name:resource_name:RAR-xxx where xxx is 1-10 digits")
                        sys.exit(99)
                    return ""
        return None

    # add given error in given fileName
    def add_error(self, error, moduleName, fileName, severity):
        if self.passNumber == 2 and self.shouldLogErrors:
            self.add_error_force(error, moduleName, fileName, severity)

    def add_error_force(self, error, moduleName, fileName, severity):
        self.jsonOutput["errors"].append( self.getFailureMsg(severity, "", error, moduleName, fileName) )

    def getFailureMsg(self, severity, waived, msg, moduleName, fileName):
        message = {}
        message["severity"] = severity
        message["waived"] = waived
        message["message"] = msg
        message["moduleName"] = moduleName
        message["fileName"] = fileName
        return message

    def logMsg(self, typ, msg):
        if self.passNumber == 2:
            self.logMsgAlways(typ, msg)

    def logMsgAlways(self, typ, msg):
        if typ == "error":
            logging.error(msg)
        elif typ == "warning":
            logging.warning(msg)
        elif typ == "info":
            logging.info(msg)
        elif typ == "debug":
            logging.debug(msg)

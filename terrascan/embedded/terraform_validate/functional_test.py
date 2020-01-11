import re
import os
import sys
import logging
import terraform_validate as t

if sys.version_info < (2, 7):
    import unittest2 as unittest
else:
    import unittest


class TestValidatorFunctional(unittest.TestCase):
    jsonOutput = {
        "failures": [],
        "errors": []
    }

    def setUp(self):
        logging.basicConfig(level=logging.CRITICAL)
        self.path = os.path.join(os.path.dirname(os.path.realpath(__file__)))
        self.jsonOutput["failures"] = []
        self.jsonOutput["errors"] = []

    def getValidatorGenerator(self, path, errorIfPropertyMissing=False):
        preprocessor = t.PreProcessor(self.jsonOutput)
        preprocessor.process(os.path.join(self.path, path))
        validator = t.Validator()
        validator.preprocessor = preprocessor
        validator.isRuleOverridden = False
        validator.overrides = []
        if errorIfPropertyMissing:
            validator.error_if_property_missing()
        return validator.get_terraform_files(False)

    def error_list_format(self, error_list):
        if type(error_list) is not list:
            error_list = [error_list]
        regex = "\n".join(map(re.escape, error_list))
        return "^{0}$".format(regex)

    #
    # the first validator line is valid; the second causes the error(s)
    #
    def test_resource_should_equal(self):
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_instance').property('value').should_equal(1)
            validator.resources('aws_instance').property('value').should_equal(2)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo.value]"):
            self.assertEqual("[aws_instance.foo.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar.value] should be '2'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[aws_instance.bar.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo.value] should be '2'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_should_equal_case_insensitive(self):
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_s3_bucket').property('value').should_equal_case_insensitive("aBc")
            validator.resources('aws_s3_bucket').property('value').should_equal("aBc")
        self.assertEqual("[aws_s3_bucket.zzz.value] should be 'aBc'. Is: 'abc'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource(self):
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').property('value').should_equal(1)
            validator.resources('aws_instance').property('nested_resource').property('value').should_equal(2)
        self.assertEqual("[aws_instance.foo.nested_resource.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"], "nested resource test failed")
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_not_equals(self):
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_instance').property('value').should_not_equal(0)
            validator.resources('aws_instance').property('value').should_not_equal(1)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo.value]"):            
            self.assertEqual("[aws_instance.foo.value] should not be '1'. Is: '1'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar.value] should not be '1'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[aws_instance.bar.value] should not be '1'. Is: '1'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo.value] should not be '1'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_not_equals_case_insensitive(self):
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_s3_bucket').property('value').should_not_equal("aBc")
            validator.resources('aws_s3_bucket').property('value').should_not_equal_case_insensitive("aBc")
        self.assertEqual("[aws_s3_bucket.zzz.value] should not be 'aBc'. Is: 'abc'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_not_equals(self):
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').property('value').should_not_equal(0)
            validator.resources('aws_instance').property('nested_resource').property('value').should_not_equal(1)
        self.assertEqual("[aws_instance.foo.nested_resource.value] should not be '1'. Is: '1'", self.jsonOutput["failures"][0]["message"], "nested resource test failed")
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_required_properties_with_list_input(self):
        required_properties_ok = ['value', 'value2']
        required_properties = ['value', 'value2', 'abc123', 'def456']
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_instance').should_have_properties(required_properties_ok)
            validator.resources('aws_instance').should_have_properties(required_properties)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo]"):
            self.assertEqual("[aws_instance.foo] should have property: 'abc123'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo] should have property: 'def456'", self.jsonOutput["failures"][1]["message"])
            self.assertEqual("[aws_instance.bar] should have property: 'abc123'", self.jsonOutput["failures"][2]["message"])
            self.assertEqual("[aws_instance.bar] should have property: 'def456'", self.jsonOutput["failures"][3]["message"])
        else:
            self.assertEqual("[aws_instance.bar] should have property: 'abc123'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar] should have property: 'def456'", self.jsonOutput["failures"][1]["message"])
            self.assertEqual("[aws_instance.foo] should have property: 'abc123'", self.jsonOutput["failures"][2]["message"])
            self.assertEqual("[aws_instance.foo] should have property: 'def456'", self.jsonOutput["failures"][3]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_excluded_properties_with_list_input(self):
        non_excluded_properties = ['value3', 'value4']
        excluded_properties = ['value', 'value2']
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_instance').should_not_have_properties(non_excluded_properties)
            validator.resources('aws_instance').should_not_have_properties(excluded_properties)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo]"):            
            self.assertEqual("[aws_instance.foo] should not have property: 'value'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo] should not have property: 'value2'", self.jsonOutput["failures"][1]["message"])
            self.assertEqual("[aws_instance.bar] should not have property: 'value'", self.jsonOutput["failures"][2]["message"])
            self.assertEqual("[aws_instance.bar] should not have property: 'value2'", self.jsonOutput["failures"][3]["message"])
        else:
            self.assertEqual("[aws_instance.bar] should not have property: 'value'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar] should not have property: 'value2'", self.jsonOutput["failures"][1]["message"])
            self.assertEqual("[aws_instance.foo] should not have property: 'value'", self.jsonOutput["failures"][2]["message"])
            self.assertEqual("[aws_instance.foo] should not have property: 'value2'", self.jsonOutput["failures"][3]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_excluded_properties_with_string_input(self):
        non_excluded_property = 'value3'
        excluded_property = 'value'
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_instance').should_not_have_properties(non_excluded_property)
            validator.resources('aws_instance').should_not_have_properties(excluded_property)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo]"):            
            self.assertEqual("[aws_instance.foo] should not have property: 'value'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar] should not have property: 'value'", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[aws_instance.bar] should not have property: 'value'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo] should not have property: 'value'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_required_properties_with_list_input(self):
        required_properties_ok = ['value', 'value2']
        required_properties = ['abc123', 'def456']
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').should_have_properties(required_properties_ok)
            validator.resources('aws_instance').property('nested_resource').should_have_properties(required_properties)
        self.assertEqual("[aws_instance.foo.nested_resource] should have property: 'abc123'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual("[aws_instance.foo.nested_resource] should have property: 'def456'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_required_properties_with_string_input(self):
        required_property_ok = 'value'
        required_property = 'def456'
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').should_have_properties(required_property_ok)
            validator.resources('aws_instance').property('nested_resource').should_have_properties(required_property)
        self.assertEqual("[aws_instance.foo.nested_resource] should have property: 'def456'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_excluded_properties_with_list_input(self):
        non_excluded_properties = ['value3', 'value4']
        excluded_properties = ['value', 'value2']
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').should_not_have_properties(non_excluded_properties)
            validator.resources('aws_instance').property('nested_resource').should_not_have_properties(excluded_properties)
        self.assertEqual("[aws_instance.foo.nested_resource] should not have property: 'value'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual("[aws_instance.foo.nested_resource] should not have property: 'value2'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_excluded_properties_with_string_input(self):
        non_excluded_property = 'value3'
        excluded_property = 'value'
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').should_not_have_properties(non_excluded_property)
            validator.resources('aws_instance').property('nested_resource').should_not_have_properties(excluded_property)
        self.assertEqual("[aws_instance.foo.nested_resource] should not have property: 'value'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_property_value_matches_regex(self):
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources('aws_instance').property('value').should_match_regex('[0-9]')
            validator.resources('aws_instance').property('value').should_match_regex('[a-z]')
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo.value]"):
            self.assertEqual("[aws_instance.foo.value] should match regex '[a-z]'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar.value] should match regex '[a-z]'", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[aws_instance.bar.value] should match regex '[a-z]'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo.value] should match regex '[a-z]'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_property_value_matches_regex(self):
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources('aws_instance').property('nested_resource').property('value').should_match_regex('[0-9]')
            validator.resources('aws_instance').property('nested_resource').property('value').should_match_regex('[a-z]')
        self.assertEqual("[aws_instance.foo.nested_resource.value] should match regex '[a-z]'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_property_invalid_json(self):
        for validator in self.getValidatorGenerator("fixtures/invalid_json"):
            validator.resources('aws_s3_bucket').property('policy').should_contain_valid_json()
        self.assertEqual("[aws_s3_bucket.invalidjson.policy] is not valid json", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_variable_substitution(self):
        for validator in self.getValidatorGenerator("fixtures/variable_substitution"):
            validator.resources('aws_instance').property('value').should_equal(1)
            validator.resources('aws_instance').property('value').should_equal(2)
        self.assertEqual("[aws_instance.foo.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_missing_variable_substitution(self):
        for validator in self.getValidatorGenerator("fixtures/missing_variable"):
            validator.resources('aws_instance').property('value').should_equal(1)
        self.assertEqual("[aws_instance.foo.value] should be '1'. Is: '!{var.missing}'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_properties_on_nonexistant_resource_type(self):
        required_properties = ['value', 'value2']
        for validator in self.getValidatorGenerator("fixtures/missing_variable"):
            validator.resources('aws_rds_instance').property('nested_resource').should_have_properties(required_properties)
        self.assertEqual(0, len(self.jsonOutput["failures"]))
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_searching_for_property_on_nonexistant_nested_resource(self):
        for validator in self.getValidatorGenerator("fixtures/resource", True):
            validator.resources('aws_instance').property('tags').property('tagname').should_equal(1)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo]"):
            self.assertEqual("[aws_instance.foo] should have property: 'tags'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar] should have property: 'tags'", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[aws_instance.bar] should have property: 'tags'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo] should have property: 'tags'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_searching_for_property_value_using_regex(self):
        for validator in self.getValidatorGenerator("fixtures/regex_variables"):
            validator.resources('aws_instance').find_property('^CPM_Service_[A-Za-z]+$').should_equal(1)
            validator.resources('aws_instance').find_property('^CPM_Service_[A-Za-z]+$').should_equal(2)
        self.assertEqual("[aws_instance.foo.CPM_Service_wibble] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_searching_for_nested_property_value_using_regex(self):
        for validator in self.getValidatorGenerator("fixtures/regex_nested_variables"):
            validator.resources('aws_instance').property('tags').find_property('^CPM_Service_[A-Za-z]+$').should_equal(1)
            validator.resources('aws_instance').property('tags').find_property('^CPM_Service_[A-Za-z]+$').should_equal(2)
        self.assertEqual("[aws_instance.foo.tags.CPM_Service_wibble] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_type_list(self):
        for validator in self.getValidatorGenerator("fixtures/resource"):
            validator.resources(['aws_instance', 'aws_elb']).property('value').should_equal(1)
            validator.resources(['aws_instance', 'aws_elb']).property('value').should_equal(2)
        if self.jsonOutput["failures"][0]["message"].startswith("[aws_instance.foo.value]"):
            self.assertEqual("[aws_instance.foo.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.bar.value] should be '2'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[aws_instance.bar.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[aws_instance.foo.value] should be '2'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual("[aws_elb.buzz.value] should be '2'. Is: '1'", self.jsonOutput["failures"][2]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_resource_type_list(self):
        for validator in self.getValidatorGenerator("fixtures/nested_resource"):
            validator.resources(['aws_instance', 'aws_elb']).property('tags').property('value').should_equal(1)
            validator.resources(['aws_instance', 'aws_elb']).property('tags').property('value').should_equal(2)
        self.assertEqual("[aws_instance.foo.tags.value] should be '2'. Is: '1'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual("[aws_elb.bar.tags.value] should be '2'. Is: '1'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_invalid_terraform_syntax(self):
        preprocessor = t.PreProcessor(self.jsonOutput)
        preprocessor.passNumber = 2
        preprocessor.shouldLogErrors = True
        preprocessor.process(os.path.join(self.path, "fixtures/invalid_syntax"))
        self.assertEqual(0, len(self.jsonOutput["failures"]))
        self.assertTrue(self.jsonOutput["errors"][0]["message"].startswith("Traceback"))

    def test_multiple_variable_substitutions(self):
        for validator in self.getValidatorGenerator("fixtures/multiple_variables"):
            validator.resources('aws_instance').property('value').should_equal(12)
            validator.resources('aws_instance').property('value').should_equal(21)
        self.assertEqual("[aws_instance.foo.value] should be '21'. Is: '12'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_nested_multiple_variable_substitutions(self):
        for validator in self.getValidatorGenerator("fixtures/multiple_variables"):
            validator.resources('aws_instance').property('value_block').property('value').should_equal(21)
            validator.resources('aws_instance').property('value_block').property('value').should_equal(12)
        self.assertEqual("[aws_instance.foo.value_block.value] should be '12'. Is: '21'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_variable_expansion(self):
        for validator in self.getValidatorGenerator("fixtures/variable_expansion"):
            validator.resources('aws_instance').property('value').should_equal(1)
            validator.resources('aws_instance').property('value').should_equal('${bar.var}')
        self.assertEqual("[aws_instance.foo.value] should be '${bar.var}'. Is: '1'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_name_matches_regex(self):
        for validator in self.getValidatorGenerator("fixtures/resource_name"):
            validator.resources('aws_foo').name_should_match_regex('^[a-z0-9_]*$')
            validator.resources('aws_instance').name_should_match_regex('^[a-z0-9_]*$')
        self.assertEqual("[aws_instance.TEST_RESOURCE] name should match regex: '^[a-z0-9_]*$'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_resource_should_not_exist(self):
        for validator in self.getValidatorGenerator("fixtures/resource_name"):
            validator.resources('xyz').should_not_exist()
            validator.resources('aws_s3_bucket').should_not_exist()
        self.assertEqual("[aws_s3_bucket] should not exist. Found in resource named badResource", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_variable_default_value_equals(self):
        for validator in self.getValidatorGenerator("fixtures/default_variable"):
            fooValue = validator.preprocessor.modulesDict["fixtures/default_variable"]["variable"].get("foo")
            barValue = validator.preprocessor.modulesDict["fixtures/default_variable"]["variable"].get("bar")
        self.assertEqual("1", fooValue)
        self.assertIsNone(barValue)
        self.assertEqual(0, len(self.jsonOutput["failures"]))
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_no_exceptions_raised_when_no_resources_present(self):
        for validator in self.getValidatorGenerator("fixtures/no_resources"):
            validator.resources('aws_instance').property('value').should_equal(1)
        self.assertEqual(0, len(self.jsonOutput["failures"]))
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_boolean_equal(self):
        for i in range(1, 6):
            for validator in self.getValidatorGenerator("fixtures/boolean_compare"):
                validator.resources("aws_db_instance").property("storage_encrypted{0}".format(i)).should_equal(True)
                validator.resources("aws_db_instance").property("storage_encrypted{0}".format(i)).should_equal("true")
                validator.resources("aws_db_instance").property("storage_encrypted{0}".format(i)).should_equal("True")
        self.assertEqual(0, len(self.jsonOutput["failures"]))
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_list_should_contain(self):
        for validator in self.getValidatorGenerator("fixtures/list_variable"):
            validator.resources("datadog_monitor").property("tags").list_should_contain(['baz:biz'])
            validator.resources("datadog_monitor").property("tags").list_should_contain('too:biz')
        if self.jsonOutput["failures"][0]["message"].startswith("[datadog_monitor.foo.tags]"):
            self.assertEqual("[datadog_monitor.foo.tags] '['baz:biz']' should contain '['too:biz']'.", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[datadog_monitor.bar.tags] '['baz:biz', 'foo:bar']' should contain '['too:biz']'.", self.jsonOutput["failures"][1]["message"])
        else:
            self.assertEqual("[datadog_monitor.bar.tags] '['baz:biz', 'foo:bar']' should contain '['too:biz']'.", self.jsonOutput["failures"][0]["message"])
            self.assertEqual("[datadog_monitor.foo.tags] '['baz:biz']' should contain '['too:biz']'.", self.jsonOutput["failures"][1]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_list_should_not_contain(self):
        for validator in self.getValidatorGenerator("fixtures/list_variable"):
            validator.resources("datadog_monitor").property("tags").list_should_not_contain(['foo:baz'])
            validator.resources("datadog_monitor").property("tags").list_should_not_contain('foo:baz')
            validator.resources("datadog_monitor").property("tags").list_should_not_contain('foo:bar')
        self.assertEqual("[datadog_monitor.bar.tags] '['baz:biz', 'foo:bar']' should not contain '['foo:bar']'.", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_property_list_scenario(self):
        for validator in self.getValidatorGenerator("fixtures/list_property"):
            validator.resources("aws_autoscaling_group").property("tag").property('propagate_at_launch').should_equal("True")
            validator.resources("aws_autoscaling_group").property("tag").property('propagate_at_launch').should_equal(True)
        self.assertEqual(0, len(self.jsonOutput["failures"]))
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_encryption_scenario(self):
        for validator in self.getValidatorGenerator("fixtures/enforce_encrypted", True):
            validator.resources("aws_db_instance_valid").property("storage_encrypted").should_equal("True")
            validator.resources("aws_db_instance_valid").property("storage_encrypted").should_equal(True)
            validator.resources("aws_db_instance_invalid").should_have_properties("storage_encrypted")
            validator.resources("aws_db_instance_invalid").property("storage_encrypted").should_equal("True")
            validator.resources("aws_db_instance_invalid2").property("storage_encrypted")
            validator.resources("aws_instance_valid").property('ebs_block_device').property("encrypted").should_equal("True")
            validator.resources("aws_instance_valid").property('ebs_block_device').property("encrypted").should_equal(True)
            validator.resources("aws_instance_invalid").should_have_properties("encrypted")
            validator.resources("aws_instance_invalid").property('ebs_block_device').property("encrypted").should_equal("True")
            validator.resources("aws_instance_invalid2").should_have_properties("storage_encrypted")
            validator.resources("aws_instance_invalid2").property('ebs_block_device').property("encrypted")
            validator.resources("aws_ebs_volume_valid").property("encrypted").should_equal("True")
            validator.resources("aws_ebs_volume_valid").property("encrypted").should_equal(True)
            validator.resources("aws_ebs_volume_invalid").should_have_properties("encrypted")
            validator.resources("aws_ebs_volume_invalid").property("encrypted").should_equal("True")
            validator.resources("aws_ebs_volume_invalid2").should_have_properties("encrypted")
            validator.resources("aws_ebs_volume_invalid2").property("encrypted")
        self.assertEqual("[aws_db_instance_invalid.foo2.storage_encrypted] should be 'True'. Is: 'False'", self.jsonOutput["failures"][0]["message"])
        self.assertEqual("[aws_db_instance_invalid2.foo3] should have property: 'storage_encrypted'", self.jsonOutput["failures"][1]["message"])
        self.assertEqual("[aws_instance_invalid.bizz2] should have property: 'encrypted'", self.jsonOutput["failures"][2]["message"])
        self.assertEqual("[aws_instance_invalid.bizz2.ebs_block_device.encrypted] should be 'True'. Is: 'False'", self.jsonOutput["failures"][3]["message"])
        self.assertEqual("[aws_instance_invalid2.bizz3] should have property: 'storage_encrypted'", self.jsonOutput["failures"][4]["message"])
        self.assertEqual("[aws_instance_invalid2.bizz3.ebs_block_device] should have property: 'encrypted'", self.jsonOutput["failures"][5]["message"])
        self.assertEqual("[aws_ebs_volume_invalid.bar2.encrypted] should be 'True'. Is: 'False'", self.jsonOutput["failures"][6]["message"])
        self.assertEqual("[aws_ebs_volume_invalid2.bar3] should have property: 'encrypted'", self.jsonOutput["failures"][7]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_with_property(self):
        for validator in self.getValidatorGenerator("fixtures/with_property"):
            validator.resources("aws_s3_bucket").with_property("acl", "private").property("policy").should_contain_valid_json()
        self.assertEqual("[aws_s3_bucket.private_bucket.policy] is not valid json", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

    def test_with_nested_property(self):
        for validator in self.getValidatorGenerator("fixtures/with_property"):
            validator.resources("aws_s3_bucket").with_property("tags", ".*'CustomTag':.*'CustomValue'.*").property("policy").should_contain_valid_json()
        self.assertEqual("[aws_s3_bucket.tagged_bucket.policy] is not valid json", self.jsonOutput["failures"][0]["message"])
        self.assertEqual(0, len(self.jsonOutput["errors"]))

# Policies

Terrascan policies are written using the [Rego policy language](https://www.openpolicyagent.org/docs/latest/policy-language/). With each rego policy a JSON "rule" file is included which defines metadata for the policy. Policies included within Terrascan are stored in the [pkg/policies/opa/rego](https://github.com/accurics/terrascan/tree/master/pkg/policies/opa/rego) directory.

## Rule JSON file

The rule files follow this naming convention: `<cloud-provider>.<resource-type>.<rule-category>.<next-available-rule-number>.json`

Here's an example of the contents of a rule file:

``` json linenums="1"
{
    "ruleName": "unrestrictedIngressAccess",
    "rule": "unrestrictedIngressAccess.rego",
    "ruleTemplate": "unrestrictedIngressAccess",
    "ruleArgument": {
      "prefix": ""
    },
    "severity": "HIGH",
    "description": "Ensure no security groups allow ingress from 0.0.0.0/0 to ALL ports and protocols",
    "ruleReferenceId": "AWS.SecurityGroup.NetworkPortsSecurity.High.0094",
    "category": "Network Ports Security",
    "version" : "1"
}
```

| Key                 | Value                                         |
| ------------------- | --------------------------------------------- |
| ruleName            | Short name for the rule                       |
| rule                | File name of the rego policy                  |
| ruleTemplate        | Rego policy template Used for the rule        |
| ruleArgument        | Argument passed to the template               |
| ruleArgument.prefix | Used for making rego policies unique          |
| severity            | Likelihood x impact of issue                  |
| description         | Description of the issue found with this rule |
| ruleReferenceId     | Unique ID of the rule in the format `<cloud-provider>.<resource-type>.<rule-category>.<next-available-rule-number>` |
| category            | Descriptive category for this rule    |
| version             | Version number for the rule/rego      |

## AWS

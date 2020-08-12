# Policies

Terrascan policies are written using the [Rego policy language](https://www.openpolicyagent.org/docs/latest/policy-language/). With each rego policy a JSON "rule" file is included which defines metadata for the policy.

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

## AWS

# Policies

Terrascan policies are written using the [Rego policy language](https://www.openpolicyagent.org/docs/latest/policy-language/). With each rego policy a JSON "rule" file is included which defines metadata for the policy. Policies included within Terrascan are stored in the [pkg/policies/opa/rego](https://github.com/accurics/terrascan/tree/master/pkg/policies/opa/rego) directory.

## Rule JSON file

The rule files follow this naming convention: `<cloud-provider>.<resource-type>.<rule-category>.<severity>.<next-available-rule-number>.json`

Here's an example of the contents of a rule file:

``` json
{
    "name": "unrestrictedIngressAccess",
    "file": "unrestrictedIngressAccess.rego",
    "template_args": {
        "prefix": ""
    },
    "severity": "HIGH",
    "description": " It is recommended that no security group allows unrestricted ingress access",
    "reference_id": "AWS.SecurityGroup.NetworkSecurity.High.0094",
    "category": "Network Ports Security",
    "version": 2
}
```

| Key                  | Value                                         |
| -------------------- | --------------------------------------------- |
| name                 | Short name for the rule                       |
| file                 | File name of the Rego policy                  |
| template_args.prefix | Used for making rego policies unique          |
| severity             | Likelihood * impact of issue                  |
| description          | Description of the issue found with this rule |
| ruleReferenceId      | Unique ID of the rule in the format `<cloud-provider>.<resource-type>.<rule-category>.<severity>.<next-available-rule-number>` |
| category            | Descriptive category for this rule    |
| version             | Version number for the rule/rego      |

--8<-- "docs/policies/aws.md"

--8<-- "docs/policies/azure.md"

--8<-- "docs/policies/gcp.md"

--8<-- "docs/policies/k8s.md"

--8<-- "docs/policies/github.md"

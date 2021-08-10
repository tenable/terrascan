# In-file Instrumentation

Terrascan can be instrumented using special commands inside your IaC files (Terraform, K8s and dockerfile)

Today, Terrascan supports these instrumentations:

* Rule Skipping
* Resource Prioritization

## Rule Skipping
Rule skipping allows you to specify a rule that should not be applied to a particular resource.

> Note:  In-file instrumentation will skip the rule only for the resource it is defined in. The `skip_rules` parameter in the config file will skip the rule for the entire scan.

### In Terraform
Use the syntax `#ts:skip=RuleID optional_comment` inside a resource to skip the rule for that resource.

#### Example
``` HCL
resource "aws_db_instance" "PtShGgAdi4" {
  #ts:skip=AWS.RDS.DataSecurity.High.0414 Reason to skip this rule
  allocated_storage       = 20
  storage_type            = "gp2"
  engine                  = "mysql"
  engine_version          = "5.7"
  instance_class          = "db.t2.micro"
 .
 .
 .
}
```
### In Kubernetes
Use the annotation
`runterrascan.io/skip:[{\"rule\": \RuleID\", \"comment\": \"reason to skip the rule\"}] ` inside a resource to skip the rule for that resource.

#### Example
``` YAML
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ingress-demo-disallowed
  annotations:
    runterrascan.io/skip: "[{\"rule\": \"AC-K8-NS-IN-H-0020\", \"comment\": \"reason to skip the rule\"}]"
spec:
  rules:
    - host: example-host.example.com
      http:
        paths:
          - backend:
              serviceName: nginx
              servicePort: 80
```
### In Dockerfile
Use the syntax `#ts:skip=RuleID optional_comment` inside the dockerfile to skip the rule for that resource.

#### Example
``` dockerfile
FROM runatlantis/atlantis:v0.16.1
#ts:skip=AC_DOCKER_0001 skip this rule.
ENV DEFAULT_TERRASCAN_VERSION=1.5.1
RUN terrascan init
ENTRYPOINT ["/bin/bash", "entrypoint.sh"]
CMD ["server"]
```
## Resource Prioritization
Resource prioritization allows you set maximum and minimum severities for violations in a given resource. Are you configuring a very sensitive resource? Set the minimum severity to `High`, so low and medium violations will be escalated. Need to suppress all violations from a particular resource? Set the maximum severity to `None`.

For maximum severity, meaningful options are Medium, Low, and None.

For minimum severity, meaningful options are High and Medium.

### In Terraform
Use the syntax `#ts:maxseverity=SEVERITY`, or `#ts:minseverity=SEVERITY` inside a resource to skip the rule for that resource.

#### Example
``` HCL
resource "aws_db_instance" "PtShGgAdi4" {
  #ts:maxseverity=Low
  allocated_storage       = 20
  storage_type            = "gp2"
  engine                  = "mysql"
  engine_version          = "5.7"
  instance_class          = "db.t2.micro"
  .
  .
  .
}
```
### In Kubernetes
Use the annotation
`runterrascan.io/minseverity: SEVERITY`, or `runterrascan.io/maxseverity: SEVERITY` inside a resource to skip the rule for that resource.

#### Example
``` YAML
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ingress-demo-disallowed
  annotations:
    runterrascan.io/minseverity: Low
spec:
  rules:
    - host: example-host.example.com
      http:
        paths:
          - backend:
              serviceName: nginx
              servicePort: 80
```
### In Dockerfile
Use the syntax `#ts:maxseverity=SEVERITY`, or `#ts:minseverity=SEVERITY` inside a dockerfile to skip the rule for that resource.

#### Example
``` dockerfile
#ts:maxseverity=None
FROM runatlantis/atlantis:v0.16.1
ENV DEFAULT_TERRASCAN_VERSION=1.5.1
RUN terrascan init
ENTRYPOINT ["/bin/bash", "entrypoint.sh"]
CMD ["server"]
```

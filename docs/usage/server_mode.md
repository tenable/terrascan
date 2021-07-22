# Using Terrascan in Server mode

[Server mode](#run-terrascan-in-server-mode) will execute Terrascan's API server. This is useful when using Terrascan to enforce a unified set of policies and configuration in multiple parts of the software development pipelines. It also simplifies programmatically interacting with Terrascan. By default the http server listens in port 9010 and supports the following routes:

> **Note:** URL placeholders are equivalent to the command line flags in the [scan command](command_line_mode.md#list-of-options-for-scan-command)



## API Routes

### Check health of server

* `GET - /health` 

### Scan IaC File

* `POST - /v1/{iac}/{iacVersion}/{cloud}/local/file/scan`

POST Parameter: `file` - Content of the file to be scanned

Example:
```curl -i -F "file=@aws_cloudfront_distribution.tf" localhost:9010/v1/terraform/v14/aws/local/file/scan```

### Scan Remote IaC

* `POST - /v1/{iac}/{iacVersion}/{cloud}/remote/dir/scan`

## Run Terrascan in Server Mode

You can launch server mode by executing the Terrascan binary, or with a Docker container. Use the following to execute the Terrascan CLI:

``` Bash
$ terrascan server
```
Use this command to launch Terrascan server mode using Docker:

``` Bash
$ docker run --rm --name terrascan -p 9010:9010 accurics/terrascan
```

### Example of how to send a request to the Terrascan server using curl:

``` Bash
$ curl -i -F "file=@aws_cloudfront_distribution.tf" localhost:9010/v1/terraform/v14/aws/local/file/scan
HTTP/1.1 100 Continue

HTTP/1.1 200 OK
Date: Sun, 16 Aug 2020 02:45:35 GMT
Content-Type: text/plain; charset=utf-8
Transfer-Encoding: chunked

{
  "results": {
    "violations": [
      {
        "rule_name": "cloudfrontNoGeoRestriction",
        "description": "Ensure that geo restriction is enabled for your Amazon CloudFront CDN distribution to whitelist or blacklist a country in order to allow or restrict users in specific locations from accessing web application content.",
        "rule_id": "AWS.CloudFront.Network Security.Low.0568",
        "severity": "LOW",
        "category": "Network Security",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoHTTPSTraffic",
        "description": "Use encrypted connection between CloudFront and origin server",
        "rule_id": "AWS.CloudFront.EncryptionandKeyManagement.High.0407",
        "severity": "HIGH",
        "category": "Encryption and Key Management",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoHTTPSTraffic",
        "description": "Use encrypted connection between CloudFront and origin server",
        "rule_id": "AWS.CloudFront.EncryptionandKeyManagement.High.0407",
        "severity": "HIGH",
        "category": "Encryption and Key Management",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoLogging",
        "description": "Ensure that your AWS Cloudfront distributions have the Logging feature enabled in order to track all viewer requests for the content delivered through the Content Delivery Network (CDN).",
        "rule_id": "AWS.CloudFront.Logging.Medium.0567",
        "severity": "MEDIUM",
        "category": "Logging",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoSecureCiphers",
        "description": "Secure ciphers are not used in CloudFront distribution",
        "rule_id": "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
        "severity": "HIGH",
        "category": "Encryption and Key Management",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      }
    ],
    "count": {
      "low": 1,
      "medium": 1,
      "high": 3,
      "total": 5
    }
  }
}
```

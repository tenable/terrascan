# Terrascan
![CI](https://github.com/accurics/terrascan/workflows/Go%20Terrascan%20build/badge.svg)
[![codecov](https://codecov.io/gh/accurics/terrascan/branch/master/graph/badge.svg)](https://codecov.io/gh/accurics/terrascan)

Terrascan is a static code analyzer and linter for security weanesses in Infrastructure as Code (IaC).

* GitHub Repo: https://github.com/accurics/terrascan
* Documentation: https://docs.accurics.com
* Discuss: https://community.accurics.com

## Features
* 500+ Policies for security best practices
* Scanning of Terraform 12+ (HCL2)
* Support for AWS, Azure, and GCP

## Installing
Terrascan's binary for your architecture can be found on the releases page. Here's an example of how to install it:

```
FIXME: Add example
```

### Homebrew
Terrascan can be installed using Homebrew on macOS:

```
brew install terrascan
```

### Chocolatey
Terrascan can be installed on Windows using Chocolatey:

```
choco install terrascan
```

### Docker
Terrascan is also available as a Docker image and can be used as follows

	$ docker run accurics/terrascan

# Using Terrascan

To scan your code for security weaknesses you can run the following

```
    $ terrascan --location tests/infrastructure/success --vars tests/infrastructure/vars.json
```

# Documentation

To learn more about Terrascan check out the documentation https://docs.accurics.com where we include a getting started guide, Terrascan's architecture, a break down of it's commands, and how to write your own policies.

# Developing Terrascan
To learn more about developing and contributing to Terrascan refer to our (contributing guide)[CONTRIBUTING.md].


To learn more about compiling Terraform and contributing suggested changes, please refer to the contributing guide.

# License

Terrascan is licensed under the [Apache 2.0 License](LICENSE).

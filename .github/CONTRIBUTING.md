# Contributing


Here are some important resources:

  * [The Tenable Community Integrations Section](https://community.tenable.com/) is a great place to discuss working with the Tenable APIs in general.
  * Bugs? [Github Issues](https://github.com/accurics/terrascan/issues) is where to report them
  * [SECURITY.md](./SECURITY.md) outlines our process for reviewing security bugs that are reported.

## Steps to contribute

1. If one doesn't already exist, [create an issue](https://github.com/accurics/terrascan/issues/new) for the bug or feature you intend to work on.
2. Create your own fork, and check it out.
3. Write your code locally. It is preferred if you create a branch for each issue or feature you work on, though not required.
4. Please add a test for any bug fix or feature being added. (Not required, but we will love you if you do)
5. Run all test cases and add any additional test cases for code you've contributed. 
6. Once all tests have passed, commit your changes to your fork and then create a Pull Request for review. Please make sure to fill out the PR template when submitting.

### Pull Requests and Code Contributions

* All tests must pass before any PR will be merged.
* Always write a clear log message for your commits. One-line messages are fine for small changes, but bigger changes should look like this:

```
    $ git commit -m "A brief summary of the commit
    > 
    > A paragraph describing what changed and its impact."
```
### Branches

The ```master``` branch is used for the current release 
Work on future releases are done on the corresponding branch name, e.g. ```1.0```, ```2.x```, etc.


### Security Testing

We have implemented a few required security checks before allowing a merge to master. 

#### Source Code Scanning

Static Code Analysis is implemented to scan the codebase for any vulnerabilities that may exist. The code base is scanned daily at a minimum to monitor for new vulnerabilities.

#### Software Composition Analysis

Software Composition Analysis is performed to monitor third party dependencies for vulnerabilities that may exist in direct or transitive dependencies. 

#### Secret Scanning

Each commit is scanned for the presence of any value that may contain a secret. If a commit contains a secret, it will be blocked from being merged. 

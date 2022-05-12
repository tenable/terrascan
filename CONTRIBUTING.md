# Contributing

Contributions are welcome, and they are greatly appreciated!

You can contribute in many ways:

## Types of Contributions

### Report Bugs

Report bugs at [https://github.com/tenable/terrascan/issues](https://github.com/tenable/terrascan/issues).

If you are reporting a bug, please include:

* Your operating system name and version.
* Any details about your local setup that might be helpful in troubleshooting.
* Detailed steps to reproduce the bug.

### Fix Bugs

Look through the GitHub issues for bugs. Anything tagged with "bug"
and "help wanted" is open to whoever wants to implement it.

### Implement Features

Look through the GitHub issues for features. Anything tagged with "enhancement"
and "help wanted" is open to whoever wants to implement it.

### Write Documentation

Terrascan could always use more documentation, whether as part of the
official Terrascan docs, or even on the web in blog posts,
articles, videos, and such. Documentation for Terrascan is located in [tenable/runterrascan.io](https://github.com/tenable/runterrascan.io) and accessible through [runterrascan.io](www.runterrascan.io). Any PRs with changes in functionality or the CLIs user interface should include a corresponding PR for documentation updates.

### Submit Feedback

The best way to send feedback is to file an issue at [https://github.com/tenable/terrascan/issues](https://github.com/tenable/terrascan/issues).

If you are proposing a feature:

* Explain in detail how it would work.
* Keep the scope as narrow as possible, to make it easier to implement.
* Remember that this is a volunteer-driven project, and that contributions
  are welcome :)

## Get Started!

Ready to contribute? Here's how to set up `terrascan` for local development.

1. Fork the `terrascan` repo on GitHub.
2. Clone your fork locally:
```
    $ git clone git@github.com:your_name_here/terrascan.git
```
3. Create a branch for local development:
```
    $ git checkout -b name-of-your-bugfix-or-feature
```
   Now you can make your changes locally.

4. When you're done making changes, check that your changes pass linting and tests. The following commands will simulate locally all checks executed as part of Terrascan's CI pipeline:
```
    $ make cicd
```
5. Commit your changes and push your branch to GitHub::
```
    $ git add .
    $ git commit -m "Your detailed description of your changes."
    $ git push origin name-of-your-bugfix-or-feature
```
6. Submit a pull request through the GitHub website.

## Pull Request Guidelines

Before you submit a pull request, check that it meets these guidelines:

1. The pull request should include tests.
2. If the pull request adds functionality or policies, the docs should be updated.
3. Make sure all tests pass by running `make cicd`.

# Release Checklist
This file provides a guideline on steps that need to be taken to properly release a new version of Terrascan:

### Run unit & integration tests
The following command will kick off both unit tests and end-to-end tests. These tests should finish clean before a release is made:

```
make cicd
```

### Bump version in source code
Once tests look clean, edit `pkg/version/version.go` and update the version number around line 22.

### Update Changelog
Running the following command will generate new entries in `CHANGELOG.md` since _since-tag_. They will be placed in the changelog file under a heading related to _future-release_. In the example below, changelogs will be generated for any changes since the `v1.4.0` tag, and placed under a section titled *v1.5.0*.

Before running this command, you'll need to have a GitHub [personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) set in the `GITHUB_TOKEN` OS environment variable.

```
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app ferrarimarco/github-changelog-generator -u tenable -p terrascan -t $GITHUB_TOKEN -b CHANGELOG.md --since-tag v1.4.0 --future-release v1.5.0
```

Next, Review `CHANGELOG.md` to ensure there are notes for the new release

Once complete, commit changes and submit a pull request. This should be the last PR before cutting the release.

### Tag the new release with git
Once the changelog PR has been merged, pull the updated code, tag it with the new version number, and push the tag back to the repo:

(again, substitute v1.5.0 for the appropriate version being released)
```
git pull
git tag v1.5.0
git push --tags
git push upstream v1.5.0
```

This will kick off the GitHub workflow to run goreleaser to perform the release.

### Brew PR

Run the commands below to update Brew to the latest Terrascan version. If you are on macOS use `shasum -a 256` instead of `sha256sum` in below command. Release version number in below command for example should be v1.5.0

```
$ export TERRASCAN_VERSION=<release_version_number>
$ brew bump-formula-pr --no-browse --url https://github.com/tenable/terrascan/archive/${TERRASCAN_VERSION}.tar.gz --sha256 $(curl -sL https://github.com/tenable/terrascan/archive/${TERRASCAN_VERSION}.tar.gz | sha256sum | awk '{print $1}') terrascan
```

### Update helm chart and kustomize directory

Manually change the version for the terrascan container image in files `deploy/helm/values.yaml` and `deploy/kustomize/base/deployment.yaml`.

# Integrating Terrascan with Pre-commit 

## Overview 
Terrascan scan can be used as a pre-commit hook in order to automatically scan your IaC before every commit. 
For more information about pre-commit hooks see https://pre-commit.com/#intro 

___

**Requirements**

* Ensure Terrascan is properly installed  (See https://runterrascan.io/docs/getting-started/#installing-terrascan)
* Have Pre-commit package manager installed (See https://pre-commit.com/#install)
___
## Integration Method 
___
### Add config file 
1. Add file called .pre-commit-config.yaml to root of repo you wish to scan with pre-commit. It should look like this: 
```yaml
repos:
    - repo: https://github.com/accurics/terrascan
        rev: <COMMIT/VERSION>  
        hooks:
        - id: terraform-pre-commit
            args: [ '-i <IAC PROVIDER>'] #optional 
```
**Note:**
The optional args line allows you to specify the IaC provider. For example, 
```yaml
repos:
    - repo: https://github.com/accurics/terrascan
        rev: <COMMIT/VERSION>  
        hooks:
        - id: terraform-pre-commit
            args: [ '-i k8s'] 
```
will cause 
```bash
'terrascan scan -i k8s' 
```
to run and thus scan kubernetes yaml files. You may exclude the args like so: 
```yaml
repos:
    - repo: https://github.com/accurics/terrascan
        rev: <COMMIT/VERSION>  
        hooks:
        - id: terraform-pre-commit
```
which causes the default 
```bash
'terrascan scan' 
```
to be run, scanning all IaC provider types. 

___

Once you have everything installed, and add the appropriate config file to your repo, 
```bash
'terrascan scan -i <IAC PROVIDER>' 
```
everytime you attempt to commit your staged changes. You can also call the hook directly on all files using pre-commit run --all-files 




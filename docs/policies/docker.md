
### docker_from
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | MEDIUM | Ensure platform flag with FROM command is not used for Docker file | AC_DOCKER_0041 | AC_DOCKER_0041 |
| Infrastructure Security | json | MEDIUM | Ensure platform flag with FROM command is not used for Docker file | AC_DOCKER_0001 | AC_DOCKER_0001 |


### docker_add
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | LOW | Ensure use Curl or Wget instead of Add to fetch packages from remote URLs, because using Add is strongly discouraged | AC_DOCKER_0035 | AC_DOCKER_0035 |
| Infrastructure Security | json | LOW | Ensure use of COPY instead of ADD unless, running a tar file | AC_DOCKER_0032 | AC_DOCKER_0032 |


### docker_expose
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | LOW | Ensure only the ports that your application needs are used and avoid exposing ports like SSH (22) | AC_DOCKER_0026 | AC_DOCKER_0026 |
| Infrastructure Security | json | HIGH | Ensure range of ports is from 0 to 65535 | AC_DOCKER_0011 | AC_DOCKER_0011 |


### docker_dockerfile
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | LOW | Ensure that HEALTHCHECK is being used. | AC_DOCKER_0047 | AC_DOCKER_0047 |
| Infrastructure Security | json | MEDIUM | Ensure the command SHELL to override the default shell instead of the RUN command. | AC_DOCKER_0020 | AC_DOCKER_0020 |


### docker_run
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | MEDIUM | Ensure Pinned Pip Package Version | AC_DOCKER_0010 | AC_DOCKER_0010 |
| Infrastructure Security | json | MEDIUM | Ensure when installing packages with pip, the '--no-cache-dir' flag should be set to make Docker images smaller | AC_DOCKER_0031 | AC_DOCKER_0031 |
| Infrastructure Security | json | HIGH | Ensure to avoid RUN with sudo command | AC_DOCKER_0007 | AC_DOCKER_0007 |
| Infrastructure Security | json | MEDIUM | Ensure WORKDIR command is getting used instead of proliferating instructions like RUN  | AC_DOCKER_0018 | AC_DOCKER_0018 |
| Infrastructure Security | json | LOW | Ensure any apt-get installs don't use '--no-install-recommends' flag | AC_DOCKER_0014 | AC_DOCKER_0014 |
| Infrastructure Security | json | MEDIUM | Ensure apt is not used with RUN command for Docker file | AC_DOCKER_0002 | AC_DOCKER_0002 |
| Infrastructure Security | json | MEDIUM | Ensure apk upgrade is not being used. | AC_DOCKER_0039 | AC_DOCKER_0039 |
| Infrastructure Security | json | MEDIUM | Ensure dnf Update is not used for Docker file | AC_DOCKER_0003 | AC_DOCKER_0003 |
| Infrastructure Security | json | MEDIUM | Ensure that there is only be one CMD instruction in a Dockerfile. If you list more than one CMD then only the last CMD will take effect | AC_DOCKER_0053 | AC_DOCKER_0053 |
| Infrastructure Security | json | MEDIUM | Ensure yum install allow manual input with RUN command for Docker file | AC_DOCKER_0004 | AC_DOCKER_0004 |
| Infrastructure Security | json | MEDIUM | Ensure instruction 'RUN <package-manager> update' should always be followed by '<package-manager> install' in the same RUN statement | AC_DOCKER_0049 | AC_DOCKER_0049 |
| Infrastructure Security | json | MEDIUM | Ensure package version is specified to avoid failures | AC_DOCKER_0033 | AC_DOCKER_0033 |
| Infrastructure Security | json | MEDIUM | Ensure Cached package data should be cleaned after installation to reduce image size | AC_DOCKER_00025 | AC_DOCKER_00025 |
| Infrastructure Security | json | MEDIUM | Ensure Yum Clean All is used after Yum Install | AC_DOCKER_0009 | AC_DOCKER_0009 |
| Infrastructure Security | json | HIGH | Ensure Yum update is being used | AC_DOCKER_0029 | AC_DOCKER_0029 |
| Infrastructure Security | json | MEDIUM | Ensure root with RUN command is not used for Docker file | AC_DOCKER_0005 | AC_DOCKER_0005 |
| Infrastructure Security | json | HIGH | Ensure that Commands 'apt-get upgrade' and 'apt-get dist-upgrade' are not being used | AC_DOCKER_0052 | AC_DOCKER_0052 |


### docker_workdir
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | HIGH | Ensure the use absolute paths for your WORKDIR. | AC_DOCKER_0013 | AC_DOCKER_0013 |


### docker_copy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | MEDIUM | Ensure not to use --chown flag when user only needs execution permission | AC_DOCKER_00024 | AC_DOCKER_00024 |



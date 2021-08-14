
### docker_from
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | MEDIUM | Ensure platform flag with FROM command is not used for Docker file | AC_DOCKER_0001 |


### docker_expose
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | HIGH | Ensure range of ports is from 0 to 65535 | AC_DOCKER_0011 |


### docker_run
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | MEDIUM | Ensure Pinned Pip Package Version | AC_DOCKER_0010 |
| Infrastructure Security | json | HIGH | Ensure to avoid RUN with sudo command | AC_DOCKER_0007 |
| Infrastructure Security | json | MEDIUM | Ensure apt is not used with RUN command for Docker file | AC_DOCKER_0002 |
| Infrastructure Security | json | MEDIUM | Ensure dnf Update is not used for Docker file | AC_DOCKER_0003 |
| Infrastructure Security | json | MEDIUM | Ensure yum install allow manual input with RUN command for Docker file | AC_DOCKER_0004 |
| Infrastructure Security | json | MEDIUM | Ensure Yum Clean All is used after Yum Install | AC_DOCKER_0009 |
| Infrastructure Security | json | MEDIUM | Ensure root with RUN command is not used for Docker file | AC_DOCKER_0005 |


### docker_workdir
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | HIGH | Ensure the use absolute paths for your WORKDIR. | AC_DOCKER_0013 |



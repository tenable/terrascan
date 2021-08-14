
### kubernetes_endpoint_slice
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | LOW | Ensure endpoint slice is not created or updated with loopback addresses as this acts as an attack vector for exploiting CVE-2021-25737 by an authorized user | AC_K8S_0113 |


### kubernetes_service
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | LOW | Ensure the use of selector is enforced for Kubernetes Ingress or LoadBalancer service | AC_K8S_0114 |
| Infrastructure Security | json | MEDIUM | Restrict the use of externalIPs | AC-K8-NS-SE-M-0188 |
| Infrastructure Security | json | MEDIUM | Ensure that the Tiller Service (Helm v2) is deleted | AC-K8-NS-SE-M-0185 |
| Infrastructure Security | json | LOW | Nodeport service can expose the worker nodes as they have public interface | AC-K8-NS-SV-L-0132 |
| Infrastructure Security | json | MEDIUM | Vulnerable to CVE-2020-8554 | AC-K8-NS-SE-M-0188 |


### kubernetes_ingress
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | MEDIUM | TLS disabled can affect the confidentiality of the data in transit | AC-K8-NS-IN-H-0020 |


### kubernetes_pod
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | MEDIUM | Containers Should Not Share the Host Network Namespace | AC-K8-NS-PO-M-0164 |
| Infrastructure Security | json | MEDIUM | Image without digest affects the integrity principle of image security | AC-K8-NS-PO-M-0133 |
| Identity and Access Management | json | HIGH | Minimize Admission of Root Containers | AC-K8-IA-PO-H-0168 |
| Security Best Practices | json | Medium | CPU Request Not Set in config file. | AC-K8-OE-PK-M-0155 |
| Security Best Practices | json | HIGH | Default Namespace Should Not be Used | AC-K8-OE-PO-M-0166 |
| Infrastructure Security | json | MEDIUM | Do Not Use CAP_SYS_ADMIN Linux Capability | AC-K8-NS-PO-H-0170 |
| Security Best Practices | json | Medium | Memory Limits Not Set in config file. | AC-K8-OE-PK-M-0158 |
| Data Protection | json | MEDIUM | Ensure That Tiller (Helm V2) Is Not Deployed | AC-K8-DS-PO-M-0177 |
| Security Best Practices | json | LOW | No readiness probe will affect automatic recovery in case of unexpected errors | AC-K8-OE-PO-L-0130 |
| Identity and Access Management | json | MEDIUM | Default seccomp profile not enabled will make the container to make non-essential system calls | AC-K8-IA-PO-M-0141 |
| Identity and Access Management | json | MEDIUM | Container images with readOnlyRootFileSystem set as false mounts the container root file system with write permissions | AC-K8-IA-PO-M-0140 |
| Infrastructure Security | json | HIGH | Prefer using secrets as files over secrets as environment variables | AC-K8-NS-PO-H-0117 |
| Infrastructure Security | json | MEDIUM | Containers Should Not Share Host IPC Namespace | AC-K8-NS-PO-M-0163 |
| Infrastructure Security | json | MEDIUM | Apply Security Context to Your Pods and Containers | AC-K8-NS-PO-M-0122 |
| Data Protection | json | MEDIUM | Ensure Kubernetes Dashboard Is Not Deployed | AC-K8-DS-PO-M-0176 |
| Identity and Access Management | json | HIGH | Allowing hostPaths to mount to Pod arise the probability of getting access to the node's filesystem | AC-K8-IA-PO-H-0138 |
| Identity and Access Management | json | MEDIUM | Some volume types mount the host file system paths to the pod or container, thus increasing the chance of escaping the container to access the host | AC-K8-IA-PO-M-0143 |
| Identity and Access Management | json | MEDIUM | Allowing the pod to make system level calls provide access to host/node sensitive information | AC-K8-IA-PO-H-0137 |
| Data Protection | json | MEDIUM | Vulnerable to CVE-2020-8555 (affected version of kube-controller-manager: v1.18.0, v1.17.0 - v1.17.4, v1.16.0 - v1.16.8, and v1.15.11 | AC-K8-DS-PO-M-0143 |
| Compliance Validation | json | MEDIUM | AlwaysPullImages plugin is not set | AC-K8-OE-PK-M-0034 |
| Identity and Access Management | json | MEDIUM | Unmasking the procMount will allow more information than is necessary to the program running in the containers spawned by k8s | AC-K8-IA-PO-M-0139 |
| Identity and Access Management | json | MEDIUM | AppArmor profile not set to default or custom profile will make the container vulnerable to kernel level threats | AC-K8-IA-PO-M-0135 |
| Identity and Access Management | json | MEDIUM | Containers Should Not Share Host Process ID Namespace | AC-K8-IA-PO-M-0162 |
| Infrastructure Security | json | MEDIUM | Containers Should Run as a High UID to Avoid Host Conflict | AC-K8-NS-PO-M-0182 |
| Identity and Access Management | json | MEDIUM | Minimize the admission of containers with the NET_RAW capability | AC-K8-IA-PS-M-0112 |
| Security Best Practices | json | LOW | No liveness probe will ensure there is no recovery in case of unexpected errors | AC-K8-OE-PO-L-0129 |
| Security Best Practices | json | LOW | No tag or container image with :Latest tag makes difficult to rollback and track | AC-K8-OE-PO-L-0134 |
| Security Best Practices | json | Medium | Memory Request Not Set in config file. | AC-K8-OE-PK-M-0157 |
| Compliance Validation | json | HIGH | Containers Should Not Run with AllowPrivilegeEscalation | AC-K8-CA-PO-H-0165 |
| Identity and Access Management | json | HIGH | Minimize the admission of privileged containers | AC-K8-IA-PO-H-0106 |
| Security Best Practices | json | Medium | CPU Limits Not Set in config file. | AC-K8-OE-PK-M-0156 |
| Infrastructure Security | json | MEDIUM | Restrict Mounting Docker Socket in a Container | AC-K8-NS-PO-M-0171 |
| Identity and Access Management | json | MEDIUM | Ensure that Service Account Tokens are only mounted where necessary | AC-K8-IA-PO-M-0105 |


### kubernetes_role
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | json | HIGH | Ensure that default service accounts are not actively used in Kubernetes Role | AC-K8-IA-RO-H-0104 |


### kubernetes_namespace
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Security Best Practices | json | LOW | No owner for namespace affects the operations | AC-K8-OE-NS-L-0128 |



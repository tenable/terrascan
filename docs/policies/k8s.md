
### kubernetes_service
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | json | MEDIUM | Restrict the use of externalIPs | AC-K8-NS-SE-M-0188 |
| Network Security | json | MEDIUM | Ensure that the Tiller Service (Helm v2) is deleted | AC-K8-NS-SE-M-0185 |
| Network Security | json | LOW | Nodeport service can expose the worker nodes as they have public interface | AC-K8-NS-SV-L-0132 |
| Network Security | json | MEDIUM | Vulnerable to CVE-2020-8554 | AC-K8-NS-SE-M-0188 |


### kubernetes_ingress
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | json | HIGH | TLS disabled can affect the confidentiality of the data in transit | AC-K8-NS-IN-H-0020 |


### kubernetes_pod
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | json | MEDIUM | Containers Should Not Share the Host Network Namespace | AC-K8-NS-PO-M-0164 |
| Network Security | json | MEDIUM | Image without digest affects the integrity principle of image security | AC-K8-NS-PO-M-0133 |
| Identity and Access Management | json | HIGH | Minimize Admission of Root Containers | AC-K8-IA-PO-H-0168 |
| Operational Efficiency | json | Medium | CPU Request Not Set in config file. | AC-K8-OE-PK-M-0155 |
| Operational Efficiency | json | MEDIUM | Default Namespace Should Not be Used | AC-K8-OE-PO-M-0166 |
| Network Security | json | HIGH | Do Not Use CAP_SYS_ADMIN Linux Capability | AC-K8-NS-PO-H-0170 |
| Operational Efficiency | json | Medium | Memory Limits Not Set in config file. | AC-K8-OE-PK-M-0158 |
| Data Security | json | MEDIUM | Ensure That Tiller (Helm V2) Is Not Deployed | AC-K8-DS-PO-M-0177 |
| Operational Efficiency | json | LOW | No readiness probe will affect automatic recovery in case of unexpected errors | AC-K8-OE-PO-L-0130 |
| Identity and Access Management | json | MEDIUM | Default seccomp profile not enabled will make the container to make non-essential system calls | AC-K8-IA-PO-M-0141 |
| Identity and Access Management | json | MEDIUM | Container images with readOnlyRootFileSystem set as false mounts the container root file system with write permissions | AC-K8-IA-PO-M-0140 |
| Network Security | json | HIGH | Prefer using secrets as files over secrets as environment variables | AC-K8-NS-PO-H-0117 |
| Network Security | json | MEDIUM | Containers Should Not Share Host IPC Namespace | AC-K8-NS-PO-M-0163 |
| Network Security | json | MEDIUM | Apply Security Context to Your Pods and Containers | AC-K8-NS-PO-M-0122 |
| Data Security | json | MEDIUM | Ensure Kubernetes Dashboard Is Not Deployed | AC-K8-DS-PO-M-0176 |
| Identity and Access Management | json | HIGH | Allowing hostPaths to mount to Pod arise the probability of getting access to the node's filesystem | AC-K8-IA-PO-H-0138 |
| Identity and Access Management | json | MEDIUM | Some volume types mount the host file system paths to the pod or container, thus increasing the chance of escaping the container to access the host | AC-K8-IA-PO-M-0143 |
| Identity and Access Management | json | HIGH | Allowing the pod to make system level calls provide access to host/node sensitive information | AC-K8-IA-PO-H-0137 |
| Operational Efficiency | json | MEDIUM | AlwaysPullImages plugin is not set | AC-K8-OE-PK-M-0034 |
| Identity and Access Management | json | MEDIUM | Unmasking the procMount will allow more information than is necessary to the program running in the containers spawned by k8s | AC-K8-IA-PO-M-0139 |
| Identity and Access Management | json | MEDIUM | AppArmor profile not set to default or custom profile will make the container vulnerable to kernel level threats | AC-K8-IA-PO-M-0135 |
| Identity and Access Management | json | MEDIUM | Containers Should Not Share Host Process ID Namespace | AC-K8-IA-PO-M-0162 |
| Network Security | json | MEDIUM | Containers Should Run as a High UID to Avoid Host Conflict | AC-K8-NS-PO-M-0182 |
| Identity and Access Management | json | MEDIUM | Minimize the admission of containers with the NET_RAW capability | AC-K8-IA-PS-M-0112 |
| Operational Efficiency | json | LOW | No liveness probe will ensure there is no recovery in case of unexpected errors | AC-K8-OE-PO-L-0129 |
| Operational Efficiency | json | LOW | No tag or container image with :Latest tag makes difficult to rollback and track | AC-K8-OE-PO-L-0134 |
| Operational Efficiency | json | Medium | Memory Request Not Set in config file. | AC-K8-OE-PK-M-0157 |
| Cloud Assets Management | json | HIGH | Containers Should Not Run with AllowPrivilegeEscalation | AC-K8-CA-PO-H-0165 |
| Identity and Access Management | json | HIGH | Minimize the admission of privileged containers | AC-K8-IA-PO-H-0106 |
| Operational Efficiency | json | Medium | CPU Limits Not Set in config file. | AC-K8-OE-PK-M-0156 |
| Network Security | json | MEDIUM | Restrict Mounting Docker Socket in a Container | AC-K8-NS-PO-M-0171 |
| Identity and Access Management | json | MEDIUM | Ensure that Service Account Tokens are only mounted where necessary | AC-K8-IA-PO-M-0105 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.120 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.116 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.117 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.106 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.110 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.111 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.107 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.112 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.108 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.109 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.105 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.113 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.118 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.114 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.115 |
| Identity and Access Management | kubernetes | MEDIUM | Container does not have resource limitations defined | accurics.kubernetes.IAM.119 |
| Data Security | json | MEDIUM | Vulnerable to CVE-2020-8555 (affected version of kube-controller-manager: v1.18.0, v1.17.0 - v1.17.4, v1.16.0 - v1.16.8,< v1.15.11 | AC-K8-DS-PO-M-0143 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.64 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.72 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.68 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.69 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.65 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.58 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.62 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.63 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.59 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.60 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.61 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.57 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.70 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.66 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.67 |
| Encryption and Key Management | kubernetes | HIGH | Container uses secrets in environment variables | accurics.kubernetes.EKM.71 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.81 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.78 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.74 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.75 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.79 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.80 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.87 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.86 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.73 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.85 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.84 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.88 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.83 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.76 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.77 |
| Identity and Access Management | kubernetes | MEDIUM | Pod has extra capabilities allowed | accurics.kubernetes.IAM.82 |


### kubernetes_role
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | json | HIGH | Ensure that default service accounts are not actively used | AC-K8-IA-RO-H-0104 |


### kubernetes_namespace
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Operational Efficiency | kubernetes | LOW | The default namespace should not be used | accurics.kubernetes.OPS.462 |
| Operational Efficiency | kubernetes | LOW | The default namespace should not be used | accurics.kubernetes.OPS.460 |
| Operational Efficiency | kubernetes | LOW | The default namespace should not be used | accurics.kubernetes.OPS.461 |
| Operational Efficiency | json | LOW | No owner for namespace affects the operations | AC-K8-OE-NS-L-0128 |



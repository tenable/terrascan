/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package testdata

var (
	// PodJSONTemplate ..
	PodJSONTemplate = []byte(`{
		"kind": "Pod",
		"apiVersion": "v1",
		"metadata": {
		  "name": "simple"
		},
		"spec": {
		  "containers": [
			{
			  "name": "healthz",
			  "image": "k8s.gcr.io/exechealthz-amd64:1.2",
			  "args": [
				"-cmd=nslookup localhost"
			  ],
			  "ports": [
				{
				  "containerPort": 8080,
				  "protocol": "TCP"
				}
			  ]
			}
		  ]
		}
	  }`)
	// PodYAMLTemplate ..
	PodYAMLTemplate = []byte(`apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  annotations:
spec:
  initContainers:
  - name: myapp-container
    image: busybox
  containers:
  - name: myapp-container
    image: nginx`)
	// CronJobYAMLTemplate ..
	CronJobYAMLTemplate = []byte(`apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure`)
	// JobYAMLTemplate ..
	JobYAMLTemplate = []byte(`apiVersion: batch/v1
kind: Job
metadata:
  name: job-wq-1
spec:
  completions: 8
  parallelism: 2
  template:
    metadata:
      name: job-wq-1
    spec:
      containers:
      - name: c
        image: gcr.io/terrascan/job-wq-1
        env:
        - name: BROKER_URL
          value: amqp://guest:guest@rabbitmq-service:5672
        - name: QUEUE
          value: job1
      restartPolicy: OnFailure`)
	// DeploymentYAMLTemplate ..
	DeploymentYAMLTemplate = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2 # tells deployment to run 2 pods matching the template
  template:
    metadata:
      labels:
        app: nginx
    spec:
      initContainers:
      - name: init
        image: busybox
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`)
	// DaemonSetYAMLTemplate ..
	DaemonSetYAMLTemplate = []byte(`apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd-elasticsearch
  namespace: kube-system
  labels:
    k8s-app: fluentd-logging
spec:
  selector:
    matchLabels:
      name: fluentd-elasticsearch
  template:
    metadata:
      labels:
        name: fluentd-elasticsearch
    spec:
      tolerations:
      # this toleration is to have the daemonset runnable on master nodes
      # remove it if your masters can't run pods
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: NoSchedule
      containers:
      - name: fluentd-elasticsearch
        image: quay.io/fluentd_elasticsearch/fluentd:v2.5.2
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      terminationGracePeriodSeconds: 30
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers`)
	// ReplicaSetYAMLTemplate ..
	ReplicaSetYAMLTemplate = []byte(`apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: frontend
  labels:
    app: guestbook
    tier: frontend
spec:
  # modify replicas according to your case
  replicas: 3
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: gcr.io/google_samples/gb-frontend:v3
`)
	// ReplicationControllerTemplate ..
	ReplicationControllerTemplate = []byte(`apiVersion: v1
kind: ReplicationController
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    app: nginx
  template:
    metadata:
      name: nginx
      labels:
        app: nginx
    spec:
      initContainers:
      - name: init1
        image: init-image-1
      - name: init2
        image: init-image-2
      - name: init3
        image: init-image-3
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
      - name: sidecar1
        image: sidecar-image-1
      - name: sidecar2
        image: sidecar-image-2

`)

	// StatefulSetTemplate ...
	StatefulSetTemplate = []byte(`apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: k8s.gcr.io/nginx-slim:0.8
        ports:
        - containerPort: 80
          name: web`)
)

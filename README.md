# hello-container

Examples and tutorials for Kubernetes use the `nginx` container image to demonstrate things like deployments, services, etc. Every time I've used a static image like `nginx` as part of a demo, I've found it challenging to answer some common questions:

* Which pod served that request?
* Which node was it running on?
* What version generated the response?

This project provides a simple "Hello World!" container that answers these questions.

## Container Image

The container image is available on Docker Hub at: [billglover/hello](https://cloud.docker.com/repository/docker/billglover/hello)

## Deployment

The sample Kubernetes deployment definition below will create a deployment with 3 replicas.

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
  labels:
    app: hello
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
      - name: hello
        image: billglover/hello
        ports:
        - containerPort: 8080
```

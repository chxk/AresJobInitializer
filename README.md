# AresJobInitializer
This repository is the implementation of initializer for [AresOperator](https://github.com/AresSys/AresOperator). It is used as initContainers of an AresJob, aiming to provide a general mechanism for startup dependency between roles.

### Prerequisites
- Docker

### Build a docker image
```shell
docker build -t ${IMAGE_NAME} .
```

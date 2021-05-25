#!/bin/bash

# install kind
sudo curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.10.0/kind-linux-amd64
sudo chmod 755 ./kind
sudo mv ./kind /usr/local/bin
sudo chown -R $USER /usr/local/bin/kind
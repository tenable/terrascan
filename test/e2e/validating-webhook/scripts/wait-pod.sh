#!/bin/bash

kubectl wait --for=condition=Ready pod/myapp --timeout=60s
#!/bin/bash

appLogs=$(kubectl logs myapp)
echo $appLogs
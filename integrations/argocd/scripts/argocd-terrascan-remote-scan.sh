#!/bin/sh
#
#    Copyright (C) 2021 Accurics, Inc.
#
#	Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#		http://www.apache.org/licenses/LICENSE-2.0
#
#	Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
###

set -o errexit
set -o nounset
set -o pipefail

TERRASCAN_SERVER="https://${SERVICE_NAME}"
IAC=${IAC_TYPE:-"k8s"}
IAC_VERSION=${IAC_VERSION:-"v1"}
CLOUD_PROVIDER=${CLOUD_PROVIDER:-"all"}
REMOTE_TYPE=${REMOTE_TYPE:-"git"}

SCAN_URL="${TERRASCAN_SERVER}/v1/${IAC}/${IAC_VERSION}/${CLOUD_PROVIDER}/remote/dir/scan"

RESPONSE=$(curl -s -w \\n%{http_code} --location -k  --request POST "$SCAN_URL" \
--header 'Content-Type: application/json' \
--data-raw '{
 "remote_type":"'${REMOTE_TYPE}'",
 "remote_url":"'${REMOTE_URL}'"
}')

echo "$RESPONSE"

# get http status code from response
HTTP_STATUS=$(printf '%s\n' "$RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" -eq 403 ]; then
    exit 3
elif [ "$HTTP_STATUS" -eq 200 ]; then
    exit 0
else
    exit 1
fi

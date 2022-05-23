#   Copyright (C) 2022 Tenable, Inc.
#
#	  Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#		http://www.apache.org/licenses/LICENSE-2.0
#
#	  Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

#!/bin/bash
set -e
count=1
declare config_file
declare copy
function fetch_configfile() {
    for i in "${@:1}"
    do
        if [[ "$i" == "-c"* ]]; then
            if [[ $i =~ -c=(.+) ]]; then
                eval config_file="${BASH_REMATCH[1]}"
                copy=${@/"$i"}
            elif [[ $i =~ -c(.+) ]]; then
                echo "unacceptable argument : $i"
                exit 1
            else
                eval var='$'$(( count + 1 ))
                eval config_file="$var"
                copy=$(echo "$@" | sed "s/ -c//")
                copy=${copy/$config_file}
            fi
        fi
    (( count += 1 ))
    done
}

fetch_configfile "$@"
if [[ ! -z $config_file ]]; then
    export TERRASCAN_CONFIG=$config_file
fi

if [[ -z $copy ]]; then
    launch-atlantis.sh $@
else
    launch-atlantis.sh $copy
fi

#!/usr/bin/env bash
set -eo pipefail
PARAMS="scan "

main() {
  initialize_
  parse_cmdline_ "$@"

  # propagate $FILES to custom function
  terrascan_ "$ARGS" "$FILES"
}

terrascan_() {
  # consume modified files passed from pre-commit so that
  # terrascan runs against only those relevant directories
  for file_with_path in $FILES; do
    file_with_path="${file_with_path// /__REPLACED__SPACE__}"
    paths[index]=$(dirname "$file_with_path")

    let "index+=1"
  done
  #put arguments array into runnable string 
  for i in "${ARGS[@]}"
  do
   PARAMS="${PARAMS} ${i}"
  done
  echo $PARAMS
  for path_uniq in $(echo "${paths[*]}" | tr ' ' '\n' | sort -u); do
    path_uniq="${path_uniq//__REPLACED__SPACE__/ }"
    pushd "$path_uniq" > /dev/null
    terrascan $PARAMS
    popd > /dev/null
  done
}

initialize_() {
  # get directory containing this script
  local dir
  local source
  source="${BASH_SOURCE[0]}"
  while [[ -L $source ]]; do # resolve $source until the file is no longer a symlink
    dir="$(cd -P "$(dirname "$source")" > /dev/null && pwd)"
    source="$(readlink "$source")"
    # if $source was a relative symlink, we need to resolve it relative to the path where the symlink file was located
    [[ $source != /* ]] && source="$dir/$source"
  done
  _SCRIPT_DIR="$(dirname "$source")"

  # source getopt function
  # shellcheck source=lib_getopt
  . "$_SCRIPT_DIR/lib_getopt"
}

parse_cmdline_() {
  declare argv
  argv=$(getopt -n Terrascan -o hi: --long help,iac-type: -- "$@") || return
  eval "set -- $argv"

  for argv; do
    case $1 in
      -i | --iac-type)   #add support for all scan options ?
        ARGS+=("$1")  #add flag 
        ARGS+=("$2") 
        shift 2       #shift up both args 
        ;;
      --)   
        shift
        FILES+=("$@")  #not sure what to do with this, replace with -f, -d handling? 
        break
        ;;
     *) 
        shift 
    esac
  done
}

# global arrays
declare -a ARGS=()
declare -a FILES=()

[[ ${BASH_SOURCE[0]} != "$0" ]] || main "$@"

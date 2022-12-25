#!/bin/bash

getopt() {
  # pure-getopt, a drop-in replacement for GNU getopt in pure Bash.
  # version 1.4.4
  #
  # Copyright 2012-2020 Aron Griffis <aron@scampersand.com>
  #
  # Permission is hereby granted, free of charge, to any person obtaining
  # a copy of this software and associated documentation files (the
  # "Software"), to deal in the Software without restriction, including
  # without limitation the rights to use, copy, modify, merge, publish,
  # distribute, sublicense, and/or sell copies of the Software, and to
  # permit persons to whom the Software is furnished to do so, subject to
  # the following conditions:
  #
  # The above copyright notice and this permission notice shall be included
  # in all copies or substantial portions of the Software.
  #
  # THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
  # OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  # MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
  # IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
  # CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
  # TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
  # SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

  _getopt_main() {
    # Returns one of the following statuses:
    #   0 success
    #   1 error parsing parameters
    #   2 error in getopt invocation
    #   3 internal error
    #   4 reserved for -T
    #
    # For statuses 0 and 1, generates normalized and shell-quoted
    # "options -- parameters" on stdout.

    declare parsed status
    declare short long='' name flags=''
    declare have_short=false

    # Synopsis from getopt man-page:
    #
    #   getopt optstring parameters
    #   getopt [options] [--] optstring parameters
    #   getopt [options] -o|--options optstring [options] [--] parameters
    #
    # The first form can be normalized to the third form which
    # _getopt_parse() understands. The second form can be recognized after
    # first parse when $short hasn't been set.

    if [[ -n ${GETOPT_COMPATIBLE+isset} || $1 == [^-]* ]]; then
      # Enable compatibility mode
      flags=c$flags
      # Normalize first to third synopsis form
      set -- -o "$1" -- "${@:2}"
    fi

    # First parse always uses flags=p since getopt always parses its own
    # arguments effectively in this mode.
    parsed=$(_getopt_parse getopt ahl:n:o:qQs:TuV \
      alternative,help,longoptions:,name:,options:,quiet,quiet-output,shell:,test,version \
      p "$@")
    status=$?
    if [[ $status != 0 ]]; then
      if [[ $status == 1 ]]; then
        echo "Try \`getopt --help' for more information." >&2
        # Since this is the first parse, convert status 1 to 2
        status=2
      fi
      return $status
    fi
    eval "set -- $parsed"

    while [[ $# -gt 0 ]]; do
      case $1 in
        (-a|--alternative)
          flags=a$flags ;;

        (-h|--help)
          _getopt_help
          return 2  # as does GNU getopt
          ;;

        (-l|--longoptions)
          long="$long${long:+,}$2"
          shift ;;

        (-n|--name)
          name=$2
          shift ;;

        (-o|--options)
          short=$2
          have_short=true
          shift ;;

        (-q|--quiet)
          flags=q$flags ;;

        (-Q|--quiet-output)
          flags=Q$flags ;;

        (-s|--shell)
          case $2 in
            (sh|bash)
              flags=${flags//t/} ;;
            (csh|tcsh)
              flags=t$flags ;;
            (*)
              echo 'getopt: unknown shell after -s or --shell argument' >&2
              echo "Try \`getopt --help' for more information." >&2
              return 2 ;;
          esac
          shift ;;

        (-u|--unquoted)
          flags=u$flags ;;

        (-T|--test)
          return 4 ;;

        (-V|--version)
          echo "pure-getopt 1.4.4"
          return 0 ;;

        (--)
          shift
          break ;;
      esac

      shift
    done

    if ! $have_short; then
      # $short was declared but never set, not even to an empty string.
      # This implies the second form in the synopsis.
      if [[ $# == 0 ]]; then
        echo 'getopt: missing optstring argument' >&2
        echo "Try \`getopt --help' for more information." >&2
        return 2
      fi
      short=$1
      have_short=true
      shift
    fi

    if [[ $short == -* ]]; then
      # Leading dash means generate output in place rather than reordering,
      # unless we're already in compatibility mode.
      [[ $flags == *c* ]] || flags=i$flags
      short=${short#?}
    elif [[ $short == +* ]]; then
      # Leading plus means POSIXLY_CORRECT, unless we're already in
      # compatibility mode.
      [[ $flags == *c* ]] || flags=p$flags
      short=${short#?}
    fi

    # This should fire if POSIXLY_CORRECT is in the environment, even if
    # it's an empty string.  That's the difference between :+ and +
    flags=${POSIXLY_CORRECT+p}$flags

    _getopt_parse "${name:-getopt}" "$short" "$long" "$flags" "$@"
  }

  _getopt_parse() {
    # Inner getopt parser, used for both first parse and second parse.
    # Returns 0 for success, 1 for error parsing, 3 for internal error.
    # In the case of status 1, still generates stdout with whatever could
    # be parsed.
    #
    # $flags is a string of characters with the following meanings:
    #   a - alternative parsing mode
    #   c - GETOPT_COMPATIBLE
    #   i - generate output in place rather than reordering
    #   p - POSIXLY_CORRECT
    #   q - disable error reporting
    #   Q - disable normal output
    #   t - quote for csh/tcsh
    #   u - unquoted output

    declare name="$1" short="$2" long="$3" flags="$4"
    shift 4

    # Split $long on commas, prepend double-dashes, strip colons;
    # for use with _getopt_resolve_abbrev
    declare -a longarr
    _getopt_split longarr "$long"
    longarr=( "${longarr[@]/#/--}" )
    longarr=( "${longarr[@]%:}" )
    longarr=( "${longarr[@]%:}" )

    # Parse and collect options and parameters
    declare -a opts params
    declare o alt_recycled=false error=0

    while [[ $# -gt 0 ]]; do
      case $1 in
        (--)
          params=( "${params[@]}" "${@:2}" )
          break ;;

        (--*=*)
          o=${1%%=*}
          if ! o=$(_getopt_resolve_abbrev "$o" "${longarr[@]}"); then
            error=1
          elif [[ ,"$long", == *,"${o#--}"::,* ]]; then
            opts=( "${opts[@]}" "$o" "${1#*=}" )
          elif [[ ,"$long", == *,"${o#--}":,* ]]; then
            opts=( "${opts[@]}" "$o" "${1#*=}" )
          elif [[ ,"$long", == *,"${o#--}",* ]]; then
            if $alt_recycled; then o=${o#-}; fi
            _getopt_err "$name: option '$o' doesn't allow an argument"
            error=1
          else
            echo "getopt: assertion failed (1)" >&2
            return 3
          fi
          alt_recycled=false
          ;;

        (--?*)
          o=$1
          if ! o=$(_getopt_resolve_abbrev "$o" "${longarr[@]}"); then
            error=1
          elif [[ ,"$long", == *,"${o#--}",* ]]; then
            opts=( "${opts[@]}" "$o" )
          elif [[ ,"$long", == *,"${o#--}::",* ]]; then
            opts=( "${opts[@]}" "$o" '' )
          elif [[ ,"$long", == *,"${o#--}:",* ]]; then
            if [[ $# -ge 2 ]]; then
              shift
              opts=( "${opts[@]}" "$o" "$1" )
            else
              if $alt_recycled; then o=${o#-}; fi
              _getopt_err "$name: option '$o' requires an argument"
              error=1
            fi
          else
            echo "getopt: assertion failed (2)" >&2
            return 3
          fi
          alt_recycled=false
          ;;

        (-*)
          if [[ $flags == *a* ]]; then
            # Alternative parsing mode!
            # Try to handle as a long option if any of the following apply:
            #  1. There's an equals sign in the mix -x=3 or -xy=3
            #  2. There's 2+ letters and an abbreviated long match -xy
            #  3. There's a single letter and an exact long match
            #  4. There's a single letter and no short match
            o=${1::2} # temp for testing #4
            if [[ $1 == *=* || $1 == -?? || \
                  ,$long, == *,"${1#-}"[:,]* || \
                  ,$short, != *,"${o#-}"[:,]* ]]; then
              o=$(_getopt_resolve_abbrev "${1%%=*}" "${longarr[@]}" 2>/dev/null)
              case $? in
                (0)
                  # Unambiguous match. Let the long options parser handle
                  # it, with a flag to get the right error message.
                  set -- "-$1" "${@:2}"
                  alt_recycled=true
                  continue ;;
                (1)
                  # Ambiguous match, generate error and continue.
                  _getopt_resolve_abbrev "${1%%=*}" "${longarr[@]}" >/dev/null
                  error=1
                  shift
                  continue ;;
                (2)
                  # No match, fall through to single-character check.
                  true ;;
                (*)
                  echo "getopt: assertion failed (3)" >&2
                  return 3 ;;
              esac
            fi
          fi

          o=${1::2}
          if [[ "$short" == *"${o#-}"::* ]]; then
            if [[ ${#1} -gt 2 ]]; then
              opts=( "${opts[@]}" "$o" "${1:2}" )
            else
              opts=( "${opts[@]}" "$o" '' )
            fi
          elif [[ "$short" == *"${o#-}":* ]]; then
            if [[ ${#1} -gt 2 ]]; then
              opts=( "${opts[@]}" "$o" "${1:2}" )
            elif [[ $# -ge 2 ]]; then
              shift
              opts=( "${opts[@]}" "$o" "$1" )
            else
              _getopt_err "$name: option requires an argument -- '${o#-}'"
              error=1
            fi
          elif [[ "$short" == *"${o#-}"* ]]; then
            opts=( "${opts[@]}" "$o" )
            if [[ ${#1} -gt 2 ]]; then
              set -- "$o" "-${1:2}" "${@:2}"
            fi
          else
            if [[ $flags == *a* ]]; then
              # Alternative parsing mode! Report on the entire failed
              # option. GNU includes =value but we omit it for sanity with
              # very long values.
              _getopt_err "$name: unrecognized option '${1%%=*}'"
            else
              _getopt_err "$name: invalid option -- '${o#-}'"
              if [[ ${#1} -gt 2 ]]; then
                set -- "$o" "-${1:2}" "${@:2}"
              fi
            fi
            error=1
          fi ;;

        (*)
          # GNU getopt in-place mode (leading dash on short options)
          # overrides POSIXLY_CORRECT
          if [[ $flags == *i* ]]; then
            opts=( "${opts[@]}" "$1" )
          elif [[ $flags == *p* ]]; then
            params=( "${params[@]}" "$@" )
            break
          else
            params=( "${params[@]}" "$1" )
          fi
      esac

      shift
    done

    if [[ $flags == *Q* ]]; then
      true  # generate no output
    else
      echo -n ' '
      if [[ $flags == *[cu]* ]]; then
        printf '%s -- %s' "${opts[*]}" "${params[*]}"
      else
        if [[ $flags == *t* ]]; then
          _getopt_quote_csh "${opts[@]}" -- "${params[@]}"
        else
          _getopt_quote "${opts[@]}" -- "${params[@]}"
        fi
      fi
      echo
    fi

    return $error
  }

  _getopt_err() {
    if [[ $flags != *q* ]]; then
      printf '%s\n' "$1" >&2
    fi
  }

  _getopt_resolve_abbrev() {
    # Resolves an abbreviation from a list of possibilities.
    # If the abbreviation is unambiguous, echoes the expansion on stdout
    # and returns 0.  If the abbreviation is ambiguous, prints a message on
    # stderr and returns 1. (For first parse this should convert to exit
    # status 2.)  If there is no match at all, prints a message on stderr
    # and returns 2.
    declare a q="$1"
    declare -a matches=()
    shift
    for a; do
      if [[ $q == "$a" ]]; then
        # Exact match. Squash any other partial matches.
        matches=( "$a" )
        break
      elif [[ $flags == *a* && $q == -[^-]* && $a == -"$q" ]]; then
        # Exact alternative match. Squash any other partial matches.
        matches=( "$a" )
        break
      elif [[ $a == "$q"* ]]; then
        # Abbreviated match.
        matches=( "${matches[@]}" "$a" )
      elif [[ $flags == *a* && $q == -[^-]* && $a == -"$q"* ]]; then
        # Abbreviated alternative match.
        matches=( "${matches[@]}" "${a#-}" )
      fi
    done
    case ${#matches[@]} in
      (0)
        [[ $flags == *q* ]] || \
        printf "$name: unrecognized option %s\\n" >&2 \
          "$(_getopt_quote "$q")"
        return 2 ;;
      (1)
        printf '%s' "${matches[0]}"; return 0 ;;
      (*)
        [[ $flags == *q* ]] || \
        printf "$name: option %s is ambiguous; possibilities: %s\\n" >&2 \
          "$(_getopt_quote "$q")" "$(_getopt_quote "${matches[@]}")"
        return 1 ;;
    esac
  }

  _getopt_split() {
    # Splits $2 at commas to build array specified by $1
    declare IFS=,
    eval "$1=( \$2 )"
  }

  _getopt_quote() {
    # Quotes arguments with single quotes, escaping inner single quotes
    declare s space='' q=\'
    for s; do
      printf "$space'%s'" "${s//$q/$q\\$q$q}"
      space=' '
    done
  }

  _getopt_quote_csh() {
    # Quotes arguments with single quotes, escaping inner single quotes,
    # bangs, backslashes and newlines
    declare s i c space
    for s; do
      echo -n "$space'"
      for ((i=0; i<${#s}; i++)); do
        c=${s:i:1}
        case $c in
          (\\|\'|!)
            echo -n "'\\$c'" ;;
          ($'\n')
            echo -n "\\$c" ;;
          (*)
            echo -n "$c" ;;
        esac
      done
      echo -n \'
      space=' '
    done
  }

  _getopt_help() {
    cat <<-EOT >&2
	Usage:
	 getopt <optstring> <parameters>
	 getopt [options] [--] <optstring> <parameters>
	 getopt [options] -o|--options <optstring> [options] [--] <parameters>
	Parse command options.
	Options:
	 -a, --alternative             allow long options starting with single -
	 -l, --longoptions <longopts>  the long options to be recognized
	 -n, --name <progname>         the name under which errors are reported
	 -o, --options <optstring>     the short options to be recognized
	 -q, --quiet                   disable error reporting by getopt(3)
	 -Q, --quiet-output            no normal output
	 -s, --shell <shell>           set quoting conventions to those of <shell>
	 -T, --test                    test for getopt(1) version
	 -u, --unquoted                do not quote the output
	 -h, --help     display this help and exit
	 -V, --version  output version information and exit
	For more details see getopt(1).
	EOT
  }

  _getopt_version_check() {
    if [[ -z $BASH_VERSION ]]; then
      echo "getopt: unknown version of bash might not be compatible" >&2
      return 1
    fi

    # This is a lexical comparison that should be sufficient forever.
    if [[ $BASH_VERSION < 2.05b ]]; then
      echo "getopt: bash $BASH_VERSION might not be compatible" >&2
      return 1
    fi

    return 0
  }

  _getopt_version_check
  _getopt_main "$@"
  declare status=$?
  unset -f _getopt_main _getopt_err _getopt_parse _getopt_quote \
    _getopt_quote_csh _getopt_resolve_abbrev _getopt_split _getopt_help \
    _getopt_version_check
  return $status
}

# vim:sw=2

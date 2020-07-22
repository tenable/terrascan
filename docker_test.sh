#!/bin/sh

test_terrascan_installed() {
    echo "Testing terrascan is installed"
    terrascan --help | grep "usage: terrascan"
    exit_code=$?
}

main() {
    exit_code=1
    test_terrascan_installed

    if [ $exit_code = 0 ]; then
	    echo "Tests passed"
    else
	    echo "Tests failed"
    fi

    exit $exit_code
}

main

#!/bin/bash
#
# Script:      run_cats.sh
# Description: downloads and runs the latest CF acceptance test suite
#
# To use this in anger, you will need to construct a CATS config file. An
# example is provided as cats_config.example.json. It is easiest to run
# the tests in a Docker image (see Dockerfile):
#
#   docker build . -t cats:devel
#   docker run -it -v $PWD:/app -e CONFIG=/app/cats_config.example.json cats:devel /app/run_cats.sh
#
set -eux

: ${GOPATH:=~/go}
: ${CONFIG:=/app/cats_config.json}
export GOPATH CONFIG 

go get -d github.com/cloudfoundry/cf-acceptance-tests
cd $GOPATH/src/github.com/cloudfoundry/cf-acceptance-tests
bin/update_submodules
bin/test
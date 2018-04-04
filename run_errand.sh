#!/bin/sh
#
# Script:      run_errand.sh
# Description: Runs a BOSH errand

set -eux

# required vars
: ${ENVIRONMENT:=test}
: $BOSH_CLIENT_SECRET $BOSH_CLIENT $BOSH_URL $CACERT $DEPLOYMENT $ERRAND

# load cert and tidy up on exit
certfile=/var/tmp/$$.ca.crt
trap 'rm -f $certfile' EXIT
echo "$CACERT" >$certfile

# set bosh env and run errand
bosh alias-env $ENVIRONMENT -e $BOSH_URL --ca-cert $certfile
bosh -e $ENVIRONMENT -d $DEPLOYMENT run-errand $ERRAND


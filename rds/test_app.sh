#!/usr/bin/env bash
set -eu
: ${RDS_ENDPOINT:=https://cf-test-rds.apps.devtest.onsclofo.uk}
curl -k "$RDS_ENDPOINT"

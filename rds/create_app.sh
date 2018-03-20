#!/usr/bin/env bash
set -eu
: ${RDS_APP_NAME:=cf-test-rds}
cf push "$RDS_APP_NAME"

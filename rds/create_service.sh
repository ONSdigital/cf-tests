#!/usr/bin/env bash
set -eu
: ${RDS_SERVICE_NAME:=test-psql}
cf create-service rds shared-psql "$RDS_SERVICE_NAME"

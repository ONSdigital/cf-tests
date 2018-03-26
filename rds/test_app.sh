#!/usr/bin/env bash
set -eux
: ${RDS_ENDPOINT:=https://cf-test-rds.apps.devtest.onsclofo.uk}
curl -s -k "$RDS_ENDPOINT" | grep "RDS service is OK"

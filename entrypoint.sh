#!/usr/bin/env bash

set -e
set -u
set -o pipefail

#if [ -n "${PARAMETER_STORE:-}" ]; then
#  export RESOLUCIONES_MID_V2_PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}/resoluciones_mid_v2/db/username --output text --query Parameter.Value)"
#  export RESOLUCIONES_MID_V2_PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/resoluciones_mid_v2/db/password --output text --query Parameter.Value)"
#fi

exec ./main "$@"
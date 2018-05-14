#! /usr/bin/env bash
set -e

# This scripts rebuilds the latest.sql file included in the schema package.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOTOP="$( cd "$DIR/../../../../../../../.." && pwd )"
go generate github.com/stellar/go/services/orbit/internal/db2/schema
go install github.com/stellar/go/services/orbit
dropdb horizon_schema --if-exists
createdb horizon_schema
DATABASE_URL=postgres://localhost/horizon_schema?sslmode=disable $GOTOP/bin/orbit db migrate up

DUMP_OPTS="--schema=public --no-owner --no-acl --inserts"
LATEST_PATH="$DIR/../db2/schema/latest.sql"
BLANK_PATH="$DIR/../test/scenarios/blank-orbit.sql"

pg_dump postgres://localhost/horizon_schema?sslmode=disable $DUMP_OPTS \
  | sed '/SET idle_in_transaction_session_timeout/d'  \
  | sed '/SET row_security/d' \
  > $LATEST_PATH
pg_dump postgres://localhost/horizon_schema?sslmode=disable \
  --clean --if-exists $DUMP_OPTS \
  | sed '/SET idle_in_transaction_session_timeout/d'  \
  | sed '/SET row_security/d' \
  > $BLANK_PATH

go generate github.com/stellar/go/services/orbit/internal/db2/schema
go generate github.com/stellar/go/services/orbit/internal/test
go install github.com/stellar/go/services/orbit

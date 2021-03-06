#! /usr/bin/env bash
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GOTOP="$( cd "$DIR/../../../../../../../.." && pwd )"
PACKAGES=$(find $GOTOP/src/github.com/rover/go/services/orbit/internal/test/scenarios -iname '*.rb' -not -name '_common_accounts.rb')
# PACKAGES=$(find $GOTOP/src/github.com/rover/go/services/orbit/internal/test/scenarios -iname 'kahuna.rb')

go install github.com/rover/go/services/orbit

dropdb hayashi_scenarios --if-exists
createdb hayashi_scenarios

export ROVER_CORE_DATABASE_URL="postgres://localhost/hayashi_scenarios?sslmode=disable"
export DATABASE_URL="postgres://localhost/horizon_scenarios?sslmode=disable"
export NETWORK_PASSPHRASE="Rover Testnet ; April 2018"
export ROVER_CORE_URL="http://localhost:8080"
export SKIP_CURSOR_UPDATE="true"

# run all scenarios
for i in $PACKAGES; do
  CORE_SQL="${i%.rb}-core.sql"
  HORIZON_SQL="${i%.rb}-orbit.sql"
  bundle exec scc -r $i --dump-root-db > $CORE_SQL

  # load the core scenario
  psql $ROVER_CORE_DATABASE_URL < $CORE_SQL

  # recreate orbit dbs
  dropdb horizon_scenarios --if-exists
  createdb horizon_scenarios

  # import the core data into orbit
  $GOTOP/bin/orbit db init
  $GOTOP/bin/orbit db rebase

  # write orbit data to sql file
  pg_dump $DATABASE_URL \
    --clean --if-exists --no-owner --no-acl --inserts \
    | sed '/SET idle_in_transaction_session_timeout/d' \
    | sed '/SET row_security/d' \
    > $HORIZON_SQL
done


# commit new sql files to bindata
go generate github.com/rover/go/services/orbit/internal/test/scenarios
# go test github.com/rover/go/services/orbit/internal/ingest

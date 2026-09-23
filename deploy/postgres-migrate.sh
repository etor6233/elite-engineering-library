#!/bin/sh
set -eu
for migration in /migrations/*.up.sql; do
  test -f "$migration"
  psql -X --set=ON_ERROR_STOP=1 --single-transaction --file="$migration"
done

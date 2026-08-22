#!/bin/sh
set -e
if [ -x /app/migrate ]; then
  /app/migrate -command up
fi
exec /app/api

#!/bin/sh
set -e

# Ensure reports dir is writable when mounted as a volume.
chown -R appuser:appuser /app/plagiarism/reports 2>/dev/null || true

# Drop privileges to appuser if setpriv is available; otherwise run as root.
if command -v setpriv >/dev/null 2>&1; then
  exec setpriv --reuid=10001 --regid=10001 --init-groups /app/server
fi

exec /app/server

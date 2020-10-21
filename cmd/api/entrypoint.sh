#!/bin/bash -e

exec > >(tee -a /var/log/app/entry.log|logger -t server -s 2>/dev/console) 2>&1

echo "[`date`] Running entrypoint script..."

echo "[`date`] Starting server..."
./api >> /var/log/app/server.log 2>&1

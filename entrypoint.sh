#!/bin/sh

# Decode the service account key at runtime
echo "$SERVICE_ACCOUNT_KEY" | base64 -d > /app/service-account.json

# Set the environment variable for Google credentials
export GOOGLE_APPLICATION_CREDENTIALS=/app/service-account.json

cat /app/service-account.json

# Start the Go application
./main
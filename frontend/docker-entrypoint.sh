#!/bin/sh
set -e

# Generate config.json from environment variables


# Start nginx
echo "Starting nginx..."
exec nginx -g "daemon off;"

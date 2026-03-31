#!/bin/bash
# deploy-wip.sh - Build and deploy local changes to the GeekBudget.WIP Portainer stack.
# Works with uncommitted and unpushed changes. No git push required.
#
# Usage: ./deploy-wip.sh

set -e

DATA_JSON="/data/data.json"
REPO_DIR="$(cd "$(dirname "$0")" && pwd)"

# Read credentials from data.json
PORTAINER_URL=$(python3 -c "import json; d=json.load(open('$DATA_JSON')); print(d['portainer']['url'])")
PORTAINER_USER=$(python3 -c "import json; d=json.load(open('$DATA_JSON')); print(d['portainer']['username'])")
PORTAINER_PASS=$(python3 -c "import json; d=json.load(open('$DATA_JSON')); print(d['portainer']['password'])")
STACK_ID=$(python3 -c "import json; d=json.load(open('$DATA_JSON')); print(d['deployments']['GeekBudget.WIP']['portainer_stack_id'])")
ENV=$(python3 -c "import json; d=json.load(open('$DATA_JSON')); print(json.dumps(d['deployments']['GeekBudget.WIP']['env']))")
WIP_URL=$(python3 -c "import json; d=json.load(open('$DATA_JSON')); print(d['deployments']['GeekBudget.WIP']['url'])")

# Refresh Portainer JWT
echo "==> Authenticating with Portainer..."
TOKEN=$(curl -sk -X POST "$PORTAINER_URL/api/auth" \
  -H "Content-Type: application/json" \
  -d "{\"Username\":\"$PORTAINER_USER\",\"Password\":\"$PORTAINER_PASS\"}" \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['jwt'])")

# Build context tar (excluding large/irrelevant dirs)
echo "==> Creating build context..."
tar -czf /tmp/geekbudget-build-context.tar.gz \
  --exclude='.git' \
  --exclude='frontend/node_modules' \
  --exclude='frontend/dist' \
  --exclude='app/node_modules' \
  --exclude='app/.next' \
  --exclude='backend/bin' \
  --exclude='tmp' \
  -C "$REPO_DIR" .

echo "    Context size: $(du -sh /tmp/geekbudget-build-context.tar.gz | cut -f1)"

# Build backend image
echo "==> Building backend image..."
RESULT=$(curl -sk -X POST \
  "$PORTAINER_URL/api/endpoints/3/docker/build?t=geekbudget-wip-backend:latest&dockerfile=backend/Dockerfile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/x-tar" \
  --data-binary @/tmp/geekbudget-build-context.tar.gz)
echo "$RESULT" | grep -E "Successfully|error" | tail -2

# Build frontend image
echo "==> Building frontend image..."
RESULT=$(curl -sk -X POST \
  "$PORTAINER_URL/api/endpoints/3/docker/build?t=geekbudget-wip-frontend:latest&dockerfile=frontend/Dockerfile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/x-tar" \
  --data-binary @/tmp/geekbudget-build-context.tar.gz)
echo "$RESULT" | grep -E "Successfully|error" | tail -2

# Build nginx image (separate context - nginx/ subdir only)
echo "==> Building nginx image..."
tar -czf /tmp/geekbudget-nginx-context.tar.gz -C "$REPO_DIR/nginx" .
RESULT=$(curl -sk -X POST \
  "$PORTAINER_URL/api/endpoints/3/docker/build?t=geekbudget-wip-nginx:latest" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/x-tar" \
  --data-binary @/tmp/geekbudget-nginx-context.tar.gz)
echo "$RESULT" | grep -E "Successfully|error" | tail -2

# Deploy stack using docker-compose.wip.yml
echo "==> Deploying stack to Portainer..."
COMPOSE=$(cat "$REPO_DIR/docker-compose.wip.yml")
HTTP_STATUS=$(curl -sk -X PUT \
  "$PORTAINER_URL/api/stacks/$STACK_ID?endpointId=3" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"stackFileContent\": $(echo "$COMPOSE" | python3 -c "import sys,json; print(json.dumps(sys.stdin.read()))"), \"env\": $ENV}" \
  -w "%{http_code}" -o /dev/null)

if [ "$HTTP_STATUS" = "200" ]; then
  echo "==> Deploy successful! App available at $WIP_URL"
else
  echo "==> Deploy failed with HTTP $HTTP_STATUS"
  exit 1
fi

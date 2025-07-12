#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="opensesame-postgres"
CONTAINER_NAME="opensesame-db"
PG_PORT="5432"
DOCKERFILE_PATH="docker/postgres/Dockerfile"

# environment variables
POSTGRES_DB="${POSTGRES_DB:-opensesame}"
POSTGRES_USER="${POSTGRES_USER:-opensesame_user}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-supersecret}"

docker build -t "${IMAGE_NAME}" -f "${DOCKERFILE_PATH}" --load .

if docker ps -a --format '{{.Names}}' \
     | grep -q "^${CONTAINER_NAME}$"; then
  echo "Removing existing container ${CONTAINER_NAME}…"
  docker rm -f "${CONTAINER_NAME}"
fi

docker run -d \
  --name "${CONTAINER_NAME}" \
  -p "${PG_PORT}:5432" \
  -e POSTGRES_DB="${POSTGRES_DB}" \
  -e POSTGRES_USER="${POSTGRES_USER}" \
  -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
  "${IMAGE_NAME}"

until docker exec "${CONTAINER_NAME}" \
       pg_isready -U "${POSTGRES_USER}" \
       > /dev/null 2>&1; do
  echo -n "."
  sleep 1
done

echo
echo "Postgres is ready and listening on port ${PG_PORT}."
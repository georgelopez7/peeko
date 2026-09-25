#!/usr/bin/env bash
# podman.sh - prints a DOCKER_HOST=unix://<socket> line pointing at the
# Docker-compatible API socket of the first running podman machine.
# Prints nothing and exits 1 when no running podman machine has a usable socket.

set -euo pipefail

machine=$(podman machine list --format '{{.Name}} {{.Running}}' 2>/dev/null | awk '$2 == "true" { print $1; exit }')

if [ -z "$machine" ]; then
  echo "podman.sh: no running podman machine found" >&2
  exit 1
fi

socket=$(podman machine inspect "$machine" | jq -r '.[0].ConnectionInfo.PodmanSocket.Path')

if [ -z "$socket" ] || [ "$socket" = "null" ] || [ ! -S "$socket" ]; then
  echo "podman.sh: podman machine '$machine' is running but has no usable socket" >&2
  exit 1
fi

echo "DOCKER_HOST=unix://$socket"

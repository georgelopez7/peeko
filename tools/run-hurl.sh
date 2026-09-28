#!/usr/bin/env bash
# Runs the hurl suite, adding a correctly signed webhook request.
#
# Signature = HMAC-SHA256(WEBHOOK_SECRET, "<timestamp>.<body>") in the configured
# encoding, optionally prefixed.
#
set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
HURL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/_hurl"
PAYLOAD='{"message":"signed from hurl","id":"evt_hurl_001"}'

if [[ -f .env ]]; then
	set -a
	# shellcheck disable=SC1091
	source .env
	set +a
else
	echo "> no .env found - copy .env.sample to .env to enable webhook signing" >&2
fi

SECRET="${WEBHOOK_SECRET:-}"
ENC="${WEBHOOK_SIGNATURE_ENCODING:-hex}"
PREFIX="${WEBHOOK_SIGNATURE_PREFIX:-}"
TS="$(date +%s)"

# NOTE: hurl's multiline string body sends "{{payload}}" plus one trailing
# newline, so the signature must cover that exact wire format.
if [[ -n "$SECRET" ]]; then
	if [[ "$ENC" == "base64" ]]; then
		SIG="$(printf '%s.%s\n' "$TS" "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET" -binary | openssl enc -base64 -A)"
	else
		SIG="$(printf '%s.%s\n' "$TS" "$PAYLOAD" | openssl dgst -sha256 -hmac "$SECRET" -hex | sed 's/^.*= //')"
	fi
else
	echo "> WEBHOOK_SECRET is empty - webhook test will be captured without a valid signature" >&2
	SIG="unsigned"
fi

args=(
	--test
	--variable "host=$BASE_URL"
	--variable "timestamp=$TS"
	--variable "signature=${PREFIX}${SIG}"
	--variable "payload=$PAYLOAD"
)

exec hurl "${args[@]}" "$HURL_DIR"/*.hurl

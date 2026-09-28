# (Test): POST a signed webhook - HMAC-SHA256 over "<timestamp>.<body>"
# The signature, timestamp and payload are computed in tools/run-hurl.sh from your .env
# (see .env.sample). Peeko annotates the result in the UI: valid / invalid / expired.
POST {{host}}/ping?source=hurl&test=webhook
User-Agent: peeko-hurl
X-Webhook-Signature: {{signature}}
X-Webhook-Timestamp: {{timestamp}}
Content-Type: application/json

```
{{payload}}
```

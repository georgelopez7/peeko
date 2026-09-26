# (Test): Ping the capture endpoint with a PATCH request
PATCH {{host}}/ping?source=hurl&test=patch
User-Agent: peeko-hurl
Content-Type: application/json
{
    "message": "patch from hurl"
}
HTTP 200

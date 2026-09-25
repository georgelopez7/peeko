# (Test): Ping the capture endpoint with a POST body and multi-value query params
POST {{host}}/ping?source=hurl&test=post
User-Agent: peeko-hurl
Content-Type: application/json
{
    "message": "hello from hurl"
}
HTTP 200

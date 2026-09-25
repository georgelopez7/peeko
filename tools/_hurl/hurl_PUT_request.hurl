# (Test): Ping the capture endpoint with a PUT request
PUT {{host}}/ping?source=hurl&test=put
User-Agent: peeko-hurl
Content-Type: application/json
{
    "message": "put from hurl"
}
HTTP 200

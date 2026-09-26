# (Test): Ping the capture endpoint with a DELETE request
DELETE {{host}}/ping?source=hurl&test=delete
User-Agent: peeko-hurl
HTTP 200

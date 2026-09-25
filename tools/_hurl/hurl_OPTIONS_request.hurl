# (Test): Ping the capture endpoint with an OPTIONS request
OPTIONS {{host}}/ping?source=hurl&test=options
User-Agent: peeko-hurl
HTTP 200

# (Test): Ping the capture endpoint with a HEAD request
HEAD {{host}}/ping?source=hurl&test=head
User-Agent: peeko-hurl
HTTP 200

# (Test): Ping the capture endpoint with a simple GET request
GET {{host}}/ping?source=hurl&test=get
User-Agent: peeko-hurl
HTTP 200

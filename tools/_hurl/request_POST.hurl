# (Test): Ping the capture endpoint with a POST body and multi-value query params
POST {{host}}/ping?source=hurl&test=post
User-Agent: peeko-hurl
Content-Type: application/json
{
    "message": "hello from hurl",
    "long_field": "peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko"
}
HTTP 200

# (Test): POST a large, deeply nested JSON payload to exercise the body panel
POST {{host}}/ping?source=hurl&test=post&variant=nested
User-Agent: peeko-hurl
Content-Type: application/json
{
    "order": {
        "id": "ord_20260930_001",
        "status": "processing",
        "customer": {
            "id": 8842,
            "profile": {
                "name": {
                    "first": "George",
                    "last": "Lopez",
                    "display": "George Lopez",
                    "aliases": ["glopez", "geo.lopez"]
                },
                "contact": {
                    "email": {
                        "primary": "george@example.com",
                        "secondary": "g.lopez@example.com",
                        "verified": [true, false]
                    },
                    "phone": {
                        "mobile": "+1-555-0100",
                        "home": "+1-555-0101",
                        "extensions": {"main": 101, "billing": 205, "support": {"tier1": 301, "tier2": 302}}
                    }
                },
                "address": {
                    "shipping": {
                        "street": {"line1": "100 Main St", "line2": {"unit": "4B", "floor": 4}},
                        "city": "Springfield",
                        "geo": {"lat": 39.7817, "lng": -89.6501, "elevation": {"meters": 182, "feet": {"approx": 597, "exact": 597.11}}}
                    },
                    "billing": {
                        "street": {"line1": "PO Box 9001", "line2": null},
                        "city": "Springfield",
                        "geo": {"lat": 39.7990, "lng": -89.6440, "elevation": {"meters": 179, "feet": {"approx": 587, "exact": 587.27}}}
                    }
                },
                "preferences": {
                    "notifications": {
                        "email": {"enabled": true, "frequency": "daily", "filters": {"categories": ["orders", "promotions", "system"], "severity": {"min": "info", "max": "critical"}}},
                        "sms": {"enabled": false, "frequency": "never", "filters": {"categories": [], "severity": {"min": "warn", "max": "critical"}}},
                        "push": {"enabled": true, "frequency": "realtime", "filters": {"categories": ["orders"], "severity": {"min": "info", "max": "warn"}, "quiet_hours": {"start": "22:00", "end": "07:00", "timezone": "America/Chicago", "override": {"weekends": true, "holidays": {"list": ["2026-12-25", "2027-01-01"], "region": "US"}}}}}
                    },
                    "locale": {"language": "en", "region": "US", "currency": {"code": "USD", "symbol": "$", "format": {"decimal": ".", "thousands": ",", "precision": 2}}}
                }
            }
        },
        "items": [
            {
                "sku": "SKU-001",
                "quantity": 2,
                "price": {"amount": 1999, "currency": "USD"},
                "attributes": {"color": "black", "size": "M", "variants": [{"id": 1, "stock": {"warehouse": {"east": 12, "west": 8, "central": {"reserve": 3, "active": 45}}}}]}
            },
            {
                "sku": "SKU-002",
                "quantity": 1,
                "price": {"amount": 499, "currency": "USD"},
                "attributes": {"color": "white", "size": "L", "variants": [{"id": 2, "stock": {"warehouse": {"east": 0, "west": 21, "central": {"reserve": 1, "active": 19}}}}]}
            }
        ],
        "totals": {
            "subtotal": {"amount": 4497, "currency": "USD"},
            "tax": {"rate": 0.0875, "amount": 393, "breakdown": {"state": 250, "county": 93, "city": 50, "special": {"district": {"id": 7, "amount": 0}}}},
            "grand": {"amount": 4890, "currency": "USD"}
        },
        "metadata": {
            "source": {"channel": "web", "referrer": {"url": "https://example.com/shop", "campaign": {"id": "cmp_42", "tags": ["spring", "sale", "2026"], "attribution": {"first_touch": "2026-09-28T10:15:00Z", "last_touch": "2026-09-30T09:42:11Z", "path": ["email", "landing", "cart", "checkout"]}}}},
            "flags": {"experimental": {"feature_a": true, "feature_b": false, "rollout": {"percent": 25, "buckets": [1, 5, 9, 13]}}}
        }
    }
}
HTTP 200

# (Test): POST with an extremely long header value to exercise value clamping
POST {{host}}/ping?source=hurl&test=post&variant=long-header
User-Agent: peeko-hurl
X-Long-Header: peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-peeko-
Content-Type: application/json
{
    "message": "long header test"
}
HTTP 200

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// rebuildSignature - computes the expected HMAC-SHA256 signature over the signed content in the configured encoding.
func (v *Verifier) rebuildSignature(body []byte, timestamp string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedContent(body, timestamp)))

	if v.cfg.Encoding == EncodingBase64 {
		return base64.StdEncoding.EncodeToString(mac.Sum(nil))
	}

	return hex.EncodeToString(mac.Sum(nil))
}

// signedContent - the exact string covered by the signature
// "<timestamp>.<body>", or just "<body>" when no timestamp is used.
func signedContent(body []byte, timestamp string) string {
	if timestamp == "" {
		return string(body)
	}

	return timestamp + "." + string(body)
}

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// rebuildSignature - computes the expected HMAC-SHA256 signature.
func (v *Verifier) rebuildSignature(body []byte, timestamp string) (string, error) {
	mac := hmac.New(sha256.New, v.cfg.Secret)
	mac.Write([]byte(signedContent(body, timestamp)))

	switch v.cfg.Encoding {
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
	case EncodingHex:
		return hex.EncodeToString(mac.Sum(nil)), nil
	default:
		return "", fmt.Errorf("invalid encoding: %q", v.cfg.Encoding)
	}
}

// signedContent - the exact string covered by the signature
// "<timestamp>.<body>", or just "<body>" when no timestamp is used.
func signedContent(body []byte, timestamp string) string {
	if timestamp == "" {
		return string(body)
	}

	return timestamp + "." + string(body)
}

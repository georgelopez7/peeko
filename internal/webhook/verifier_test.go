package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testBody = `{"event":"ping","id":"evt_123"}`
const testSecret = "whsec_test_secret"

var testTime = time.Unix(1700000000, 0)

// testConfig - a verifier config signing "<timestamp>.<body>".
func testConfig() Config {
	return Config{
		Secret:          []byte(testSecret),
		SignatureHeader: "X-Webhook-Signature",
		TimestampHeader: "X-Webhook-Timestamp",
		Encoding:        EncodingHex,
		Tolerance:       5 * time.Minute,
	}
}

// testHeaders - headers signing ts + "." + body.
func testHeaders(tsHeader string, ts int64, sig string) http.Header {
	h := http.Header{}
	h.Set("X-Webhook-Signature", sig)
	h.Set(tsHeader, strconv.FormatInt(ts, 10))
	return h
}

// signHex - computes a hex HMAC-SHA256 of the signed payload.
func signHex(secret string, signed string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	return hex.EncodeToString(mac.Sum(nil))
}

// signBase64 - computes a base64 HMAC-SHA256 of the signed payload.
func signBase64(secret string, signed string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func TestVerifier(t *testing.T) {
	t.Run("should be disabled without a secret", func(t *testing.T) {
		require.False(t, Config{}.Enabled())
	})

	t.Run("should be enabled with a secret", func(t *testing.T) {
		require.True(t, Config{Secret: []byte("s")}.Enabled())
	})

	t.Run("should skip when verification is disabled", func(t *testing.T) {
		v := NewVerifier(Config{})

		res := v.Verify(http.Header{}, []byte(testBody), testTime)

		require.Equal(t, StatusSkipped, res.Status)
		require.NotEmpty(t, res.Note)
	})

	t.Run("should skip when the verifier is nil", func(t *testing.T) {
		var v *Verifier

		res := v.Verify(http.Header{}, []byte(testBody), testTime)

		require.Equal(t, StatusSkipped, res.Status)
	})

	t.Run("should flag a missing signature header", func(t *testing.T) {
		v := NewVerifier(testConfig())

		h := http.Header{}
		h.Set("X-Webhook-Timestamp", "1700000000")

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusMissing, res.Status)
		require.Contains(t, res.Note, "X-Webhook-Signature")
	})

	t.Run("should flag a missing timestamp header", func(t *testing.T) {
		v := NewVerifier(testConfig())

		h := http.Header{}
		h.Set("X-Webhook-Signature", "deadbeef")

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusMissing, res.Status)
		require.Contains(t, res.Note, "X-Webhook-Timestamp")
	})

	t.Run("should accept a valid signature", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
		require.Equal(t, sig, res.ComputedSignature)
		require.Equal(t, signed, res.SignedPayload)
		require.Empty(t, res.Note)
	})

	t.Run("should accept timestamps within tolerance in the past and future", func(t *testing.T) {
		v := NewVerifier(testConfig())

		for _, ts := range []int64{testTime.Unix() - int64(4*time.Minute/time.Second), testTime.Unix() + int64(4*time.Minute/time.Second)} {
			signed := strconv.FormatInt(ts, 10) + "." + testBody
			sig := signHex(testSecret, signed)
			res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)
			require.Equal(t, StatusValid, res.Status, "ts=%d", ts)
		}
	})

	t.Run("should flag an expired timestamp", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix() - int64(10*time.Minute/time.Second)
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)

		require.Equal(t, StatusExpired, res.Status)
		require.Contains(t, res.Note, "outside tolerance")
	})

	t.Run("should flag an unparsable timestamp", func(t *testing.T) {
		v := NewVerifier(testConfig())

		signed := "not-a-number." + testBody
		sig := signHex(testSecret, signed)

		h := http.Header{}
		h.Set("X-Webhook-Signature", sig)
		h.Set("X-Webhook-Timestamp", "not-a-number")

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusExpired, res.Status)
		require.Contains(t, res.Note, "unix")
	})

	t.Run("should flag a tampered body", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(`{"event":"hacked"}`), testTime)

		require.Equal(t, StatusInvalid, res.Status)
		require.Equal(t, sig, res.ReceivedSignature)
		require.NotEqual(t, sig, res.ComputedSignature)
	})

	t.Run("should flag a tampered timestamp", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts+1, sig), []byte(testBody), testTime)

		require.Equal(t, StatusInvalid, res.Status)
	})

	t.Run("should accept a body-only signature", func(t *testing.T) {
		cfg := testConfig()

		cfg.TimestampHeader = ""
		v := NewVerifier(cfg)
		sig := signHex(testSecret, testBody)

		h := http.Header{}
		h.Set("X-Webhook-Signature", sig)

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
		require.Equal(t, testBody, res.SignedPayload)
		require.Empty(t, res.Note)
	})

	t.Run("should ignore the timestamp header in body-only mode", func(t *testing.T) {
		cfg := testConfig()

		cfg.TimestampHeader = ""
		v := NewVerifier(cfg)
		sig := signHex(testSecret, testBody)

		h := http.Header{}
		h.Set("X-Webhook-Signature", sig)
		h.Set("X-Webhook-Timestamp", strconv.FormatInt(testTime.Unix()-3600, 10))

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
	})

	t.Run("should accept a prefixed signature", func(t *testing.T) {
		cfg := testConfig()

		cfg.Prefix = "sha256="
		v := NewVerifier(cfg)
		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := "sha256=" + signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
		require.Contains(t, res.ReceivedSignature, "sha256=")
	})

	t.Run("should accept a signature without the configured prefix", func(t *testing.T) {
		cfg := testConfig()

		cfg.Prefix = "sha256="
		v := NewVerifier(cfg)
		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
	})

	t.Run("should accept a base64-encoded signature", func(t *testing.T) {
		cfg := testConfig()

		cfg.Encoding = EncodingBase64
		v := NewVerifier(cfg)
		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signBase64(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
		require.Equal(t, sig, res.ComputedSignature)
	})

	t.Run("should flag an encoding mismatch", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := signBase64(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), []byte(testBody), testTime)

		require.Equal(t, StatusInvalid, res.Status)
		require.NotEqual(t, sig, res.ComputedSignature)
	})

	t.Run("should flag a malformed signature", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, "g7h8!not-hex"), []byte(testBody), testTime)

		require.Equal(t, StatusInvalid, res.Status)
		require.Contains(t, res.Note, "mismatch")
	})

	t.Run("should flag rotation-style multi-signatures as invalid", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		oldSig := signHex("old_secret", signed)
		newSig := signHex(testSecret, signed)

		h := http.Header{}
		h.Set("X-Webhook-Signature", oldSig+" "+newSig)
		h.Set("X-Webhook-Timestamp", strconv.FormatInt(ts, 10))

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusInvalid, res.Status)
	})

	t.Run("should accept a comma-prefixed signature", func(t *testing.T) {
		cfg := testConfig()

		cfg.Prefix = "v1,"
		v := NewVerifier(cfg)
		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "." + testBody
		sig := "v1," + signHex(testSecret, signed)

		h := http.Header{}
		h.Set("X-Webhook-Signature", sig)
		h.Set("X-Webhook-Timestamp", strconv.FormatInt(ts, 10))

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
	})

	t.Run("should accept an empty body", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()
		signed := strconv.FormatInt(ts, 10) + "."
		sig := signHex(testSecret, signed)

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, sig), nil, testTime)

		require.Equal(t, StatusValid, res.Status)
	})

	t.Run("should explain an invalid signature on an empty body", func(t *testing.T) {
		v := NewVerifier(testConfig())

		ts := testTime.Unix()

		res := v.Verify(testHeaders("X-Webhook-Timestamp", ts, "00"), nil, testTime)

		require.Equal(t, StatusInvalid, res.Status)
		require.NotEmpty(t, res.Note)
	})

	t.Run("should read a multi-value signature header", func(t *testing.T) {
		cfg := testConfig()

		cfg.TimestampHeader = ""
		v := NewVerifier(cfg)
		sig := signHex(testSecret, testBody)

		h := http.Header{}
		h["X-Webhook-Signature"] = []string{sig}

		res := v.Verify(h, []byte(testBody), testTime)

		require.Equal(t, StatusValid, res.Status)
	})
}

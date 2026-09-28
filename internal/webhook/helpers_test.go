package webhook

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignedContent(t *testing.T) {
	t.Run("should join the timestamp and body with a dot", func(t *testing.T) {
		require.Equal(t, "1700000000."+testBody, signedContent([]byte(testBody), "1700000000"))
	})

	t.Run("should return only the body without a timestamp", func(t *testing.T) {
		require.Equal(t, testBody, signedContent([]byte(testBody), ""))
	})

	t.Run("should accept an empty body", func(t *testing.T) {
		require.Equal(t, "", signedContent(nil, ""))
		require.Equal(t, "1700000000.", signedContent(nil, "1700000000"))
	})
}

func TestRebuildSignature(t *testing.T) {
	body := []byte(testBody)
	ts := "1700000000"

	t.Run("should build the same hex signature the verifier accepts", func(t *testing.T) {
		v := NewVerifier(Config{Secret: []byte(testSecret)})

		built, err := v.rebuildSignature(body, ts)

		require.NoError(t, err)
		require.Equal(t, signHex(testSecret, ts+"."+testBody), built)
	})

	t.Run("should build a base64 signature when configured", func(t *testing.T) {
		cfg := testConfig()
		cfg.Encoding = EncodingBase64
		v := NewVerifier(cfg)

		built, err := v.rebuildSignature(body, ts)

		require.NoError(t, err)
		require.Equal(t, signBase64(testSecret, ts+"."+testBody), built)
	})

	t.Run("should sign only the body when no timestamp is given", func(t *testing.T) {
		v := NewVerifier(Config{Secret: []byte(testSecret)})

		built, err := v.rebuildSignature(body, "")

		require.NoError(t, err)
		require.Equal(t, signHex(testSecret, testBody), built)
	})

	t.Run("should use the configured secret", func(t *testing.T) {
		cfg := testConfig()
		v := NewVerifier(cfg)

		built, err := v.rebuildSignature(body, ts)

		require.NoError(t, err)
		require.Equal(t, signHex(testSecret, ts+"."+testBody), built)
		require.NotEqual(t, signHex("other_secret", ts+"."+testBody), built)
	})

	t.Run("should reject an invalid encoding", func(t *testing.T) {
		cfg := testConfig()
		cfg.Encoding = Encoding("md5")
		v := Verifier{cfg: cfg}

		built, err := v.rebuildSignature(body, ts)

		require.Error(t, err)
		require.Empty(t, built)
	})

	t.Run("should produce a signature a full verifier round-trip accepts", func(t *testing.T) {
		cfg := testConfig()
		v := NewVerifier(cfg)

		built, err := v.rebuildSignature(body, ts)
		require.NoError(t, err)

		h := testHeaders("X-Webhook-Timestamp", testTime.Unix(), built)
		res := v.Verify(h, body, testTime)

		require.Equal(t, StatusValid, res.Status)
	})
}

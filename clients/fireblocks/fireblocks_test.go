package fireblocks

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewClient(t *testing.T) {
	cl := NewClient("", "awesome-api-key-1234", nil)

	require.NotNil(t, cl)
}

func Test_NewClient_signJWT(t *testing.T) {
	cl := NewClient("", "", nil)

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	privBytes := x509.MarshalPKCS1PrivateKey(key)

	privPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	privPEM := pem.EncodeToMemory(privPem)
	cl.LoadPrivateKey(privPEM)

	accessToken, err := cl.signJWT("/v1/transactions", []byte(`{ "key": "value" }`))
	require.NoError(t, err)

	tokenToValidate, err := cl.verify(accessToken, PublicKeyToBytes(t, &key.PublicKey))
	require.NoError(t, err)

	// Check if the token is valid.
	require.Equal(t, true, tokenToValidate.Valid)
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(t *testing.T, pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	require.NoError(t, err)

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

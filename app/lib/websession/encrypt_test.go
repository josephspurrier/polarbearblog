package websession

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	key := "59f3726ba3f8271ddf32224b809c42e9ef4523865c74cb64e9d7d5a031f1f706"
	raw := []byte("hello")

	en := NewEncryptedStorage(key)

	enc, err := en.Encrypt(raw)
	assert.NoError(t, err)

	dec, err := en.Decrypt(enc)
	assert.NoError(t, err)

	assert.Equal(t, raw, dec)
}

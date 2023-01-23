package votifier

import (
	"crypto/rsa"
	"math/rand"
	"reflect"
	"testing"
)

// Extremely bad random number generation, used only for testing purposes.
type badRandomReader struct{}

func (badRandomReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(rand.Intn(255))
	}
	return len(p), nil
}

func TestEncodeV1(t *testing.T) {
	// Generate a set of keys for later use
	key, err := rsa.GenerateKey(new(badRandomReader), 2048)
	if err != nil {
		t.Error(err)
		return
	}

	// Try to encrypt this vote.
	v := Vote{
		ServiceName: "golang",
		Username:    "golang",
		Address:     "127.0.0.1",
	}
	s, err := v.EncodeV1(&key.PublicKey)
	if err != nil {
		t.Error(err)
		return
	}

	if len(*s) != 256 {
		t.Errorf("Encrypted PKCS1v15 output should be 256 bytes, but it is %d bytes long", len(*s))
		return
	}

	// Try to decrypt this vote.
	var d Vote
	err = d.DecodeV1(*s, key)
	if err != nil {
		t.Error(err)
		return
	}

	if reflect.DeepEqual(v, d) {
		t.Error("votes don't match")
		return
	}
}

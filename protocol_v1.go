package votifier

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strings"
)

// DecodeV1 decodes the vote from the V1 protocol.
func (v *Vote) DecodeV1(data []byte, key *rsa.PrivateKey) error {
	if v == nil {
		*v = Vote{}
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return fmt.Errorf("failed to decrypt vote: %w", err)
	}

	elements := strings.Split(string(decrypted), "\n")
	if len(elements) != 6 {
		return fmt.Errorf("invalid element count, wanted 6, got %d", len(elements))
	}
	if elements[0] != "VOTE" {
		return fmt.Errorf("first element is incorrect; expected 'VOTE', got %s", elements[0])
	}
	v.ServiceName = elements[1]
	v.Username = elements[2]
	v.Address = elements[3]
	v.Timestamp = parseTime(elements[4])
	return nil
}

// EncodeV1 encodes the vote to the V1 protocol.
func (v *Vote) EncodeV1(publicKey *rsa.PublicKey) (*[]byte, error) {
	if v.Timestamp.IsZero() {
		v.Timestamp = timeNow()
	}

	s := strings.Join([]string{
		"VOTE",
		v.ServiceName,
		v.Username,
		v.Address,
		formatTimeMillis(v.Timestamp),
		"",
	}, "\n")
	msg := []byte(s)

	// Encrypt the v using the supplied public key.
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
	if err != nil {
		return nil, err
	}

	return &encrypted, nil
}

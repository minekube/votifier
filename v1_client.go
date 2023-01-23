package votifier

import (
	"crypto/rsa"
	"fmt"
)

// V1Client represents a Votifier v1 client.
type V1Client struct {
	address   string
	publicKey *rsa.PublicKey
}

// NewV1Client creates a new Votifier client.
func NewV1Client(address string, publicKey *rsa.PublicKey) *V1Client {
	return &V1Client{address, publicKey}
}

// SendVote sends a vote through the client.
func (client *V1Client) SendVote(vote Vote) error {
	conn, err := dial(client.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	serialized, err := vote.EncodeV1(client.publicKey)
	if err != nil {
		return err
	}

	_, err = conn.Write(*serialized)
	if err != nil {
		return fmt.Errorf("failed to send vote: %w", err)
	}
	return nil
}

package votifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// V2Client represents a Votifier v2 client.
type V2Client struct {
	address string
	token   string
}

type v2Response struct {
	Status string `json:"status"`
	Cause  string `json:"cause"`
	Error  string `json:"error"`
}

// NewV2Client creates a new Votifier v2 client.
func NewV2Client(address string, token string) *V2Client {
	return &V2Client{address, token}
}

// SendVote sends a vote through the client.
func (client *V2Client) SendVote(vote Vote) error {
	conn, err := dial(client.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	greeting := make([]byte, 64)
	read, err := conn.Read(greeting)
	if err != nil {
		return fmt.Errorf("error reading greeting: %w", err)
	}

	parts := bytes.Split(greeting[:read-1], []byte(" "))
	if len(parts) != 3 {
		return errors.New("not a v2 server")
	}
	challenge := string(parts[2])

	serialized, err := vote.EncodeV2(client.token, challenge)
	if err != nil {
		return fmt.Errorf("error encoding vote: %w", err)
	}
	_, err = conn.Write(serialized)
	if err != nil {
		return fmt.Errorf("failed to send vote: %w", err)
	}

	// read response
	resBuf := make([]byte, 256)
	read, err = conn.Read(resBuf)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var res v2Response
	rd := bytes.NewReader(resBuf[:read])
	err = json.NewDecoder(rd).Decode(&res)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	if !strings.EqualFold(res.Status, "ok") {
		return fmt.Errorf("remote server error: %w", &remoteError{
			cause: res.Cause,
			err:   errors.New(res.Error),
		})
	}

	return nil
}

type remoteError struct {
	cause string
	err   error
}

func (e *remoteError) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.err)
}

func (e *remoteError) Unwrap() error {
	return e.err
}

package votifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type votifier2Wrapper struct {
	Payload   string `json:"payload"`
	Signature []byte `json:"signature"`
}

type votifier2Inner struct {
	ServiceName string `json:"serviceName"`
	Username    string `json:"username"`
	Address     string `json:"address"`
	Timestamp   int64  `json:"timestamp"`
	Challenge   string `json:"challenge"`
}

const v2Magic int16 = 0x733A

func (v *Vote) DecodeV2(data []byte, tokenProvider TokenProvider, challenge string) error {
	rd := bytes.NewReader(data)

	// verify v2 magic
	var magicRead int16
	err := binary.Read(rd, binary.BigEndian, &magicRead)
	if err != nil {
		return err
	}

	if magicRead != v2Magic {
		return errors.New("v2 magic mismatch")
	}

	// read message length
	var length int16
	if err = binary.Read(rd, binary.BigEndian, &length); err != nil {
		return err
	}

	// now for the fun part
	var wrapper votifier2Wrapper
	if err = json.NewDecoder(rd).Decode(&wrapper); err != nil {
		return err
	}

	var vote votifier2Inner
	if err = json.NewDecoder(strings.NewReader(wrapper.Payload)).Decode(&vote); err != nil {
		return err
	}

	// validate challenge
	if vote.Challenge != challenge {
		return errors.New("invalid challenge")
	}

	// validate HMAC
	token := tokenProvider.Token(vote.ServiceName)
	m := hmac.New(sha256.New, []byte(token))
	m.Write([]byte(wrapper.Payload))
	s := m.Sum(nil)
	if !hmac.Equal(s, wrapper.Signature) {
		return errors.New("invalid signature")
	}

	v.ServiceName = vote.ServiceName
	v.Username = vote.Username
	v.Address = vote.Address
	v.Timestamp = time.UnixMilli(vote.Timestamp)
	return nil
}

func (v *Vote) EncodeV2(token string, challenge string) ([]byte, error) {
	if v.Timestamp.IsZero() {
		v.Timestamp = timeNow()
	}
	inner := votifier2Inner{
		ServiceName: v.ServiceName,
		Address:     v.Address,
		Username:    v.Username,
		Timestamp:   v.Timestamp.UnixMilli(),
		Challenge:   challenge,
	}

	// encode inner vote and generate outer package
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(inner); err != nil {
		return nil, err
	}

	innerJSON := buf.String()
	m := hmac.New(sha256.New, []byte(token))
	_, err := buf.WriteTo(m)
	if err != nil {
		return nil, fmt.Errorf("failed to write to hmac: %w", err)
	}

	wrapper := votifier2Wrapper{
		Payload:   innerJSON,
		Signature: m.Sum(nil),
	}

	// assemble full package
	var wrapperBuf bytes.Buffer
	if err := json.NewEncoder(&wrapperBuf).Encode(wrapper); err != nil {
		return nil, fmt.Errorf("failed to encode wrapper: %w", err)
	}

	buf.Reset()
	if err := binary.Write(buf, binary.BigEndian, v2Magic); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, int16(wrapperBuf.Len())); err != nil {
		return nil, err
	}
	_, err = wrapperBuf.WriteTo(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}

	return buf.Bytes(), nil
}

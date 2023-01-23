package votifier

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// Protocol represents a Votifier protocol.
type Protocol int

// Protocols versions supported.
const (
	V1 Protocol = iota + 1 // Uses base64 encoded RSA public key as token for vote verification.
	V2                     // Uses any token string for vote verification.
)

// VoteListener takes a vote and an int describing the protocol version (1 or 2).
type VoteListener func(*Vote, Protocol) error

type ReceiverRecord struct {
	PrivateKey    *rsa.PrivateKey // v1
	TokenProvider TokenProvider   // v2
}

// Server represents a Votifier server.
type Server struct {
	VoteHandler VoteListener // Required vote handler
	Records     []ReceiverRecord
	OnErr       func(net.Conn, error) // Optional connection handler
}

// ListenAndServe binds to a specified address-port pair and starts serving Votifier requests.
func (s *Server) ListenAndServe(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()
	return s.Serve(l)
}

// Serve serves requests on the provided listener.
func (s *Server) Serve(ln net.Listener) error {
	for {
		// Wait for a connection.
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()
	if err := s.HandleConn(c); err != nil && s.OnErr != nil {
		s.OnErr(c, err)
	}
}

func (s *Server) HandleConn(c net.Conn) error {
	challenge, err := randomString()
	if err != nil {
		// something very bad happened - only caused when /dev/urandom
		// also returns an error, which should never happen.
		return fmt.Errorf("error generating challenge: %v", err)
	}
	err = c.SetDeadline(timeNow().Add(5 * time.Second))
	if err != nil {
		return fmt.Errorf("error setting deadline: %v", err)
	}

	// Write greeting
	_, err = fmt.Fprintf(c, "VOTIFIER 2 %s\n", challenge)
	if err != nil {
		return fmt.Errorf("error writing greeting: %v", err)
	}

	// Read in what data we can and try to handle it
	data := make([]byte, 1024)
	read, err := c.Read(data)
	if err != nil {
		return fmt.Errorf("error reading data: %v", err)
	}

	// Do we have v2 magic?
	reader := bytes.NewReader(data[:2])
	var magicRead int16
	if err = binary.Read(reader, binary.BigEndian, &magicRead); err != nil {
		return fmt.Errorf("error reading magic: %v", err)
	}

	isv2 := magicRead == v2Magic
	for _, record := range s.Records {
		v := new(Vote)
		if !isv2 && record.PrivateKey != nil {
			err = v.DecodeV1(data[:read], record.PrivateKey)
			if err != nil {
				continue
			}
			err = s.VoteHandler(v, V1)
			continue
		} else {
			err = v.DecodeV2(data[:read], record.TokenProvider, challenge)
			if err != nil {
				continue
			}

			err = s.VoteHandler(v, V2)
			if err != nil {
				continue
			}

			_, _ = io.WriteString(c, `{"status":"ok"}`)
			return nil
		}
	}

	// We couldn't decrypt it correctly
	if isv2 {
		result := v2Response{
			Status: "error",
			Cause:  "decode",
		}
		if err != nil {
			result.Error = fmt.Sprint(err)
		}
		_ = json.NewEncoder(c).Encode(result)
	}
	return err
}

type Result struct {
	Status string
}

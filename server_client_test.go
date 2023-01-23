package votifier

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
)

var (
	Protocols = []Protocol{V1, V2}
)

func TestServer(t *testing.T) {
	// Generate a set of keys for later use
	key, err := rsa.GenerateKey(new(badRandomReader), 2048)
	if err != nil {
		t.Error(err)
		return
	}

	for _, i := range Protocols {
		t.Run(fmt.Sprintf("Protocol %d", i), func(t *testing.T) {
			v := Vote{
				ServiceName: "golang",
				Username:    "golang",
				Address:     "127.0.0.1",
			}
			vl := func(rv *Vote, ver Protocol) error {
				if reflect.DeepEqual(v, *rv) {
					t.Error("Vote received did not match original")
				}

				if ver != i {
					t.Errorf("Vote is not v %d", i)
				}
				return nil
			}

			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Error(err)
			}
			defer listener.Close()

			var client Client
			switch i {
			case V1:
				pk := key.PublicKey
				client = NewV1Client(listener.Addr().String(), &pk)
			case V2:
				client = NewV2Client(listener.Addr().String(), "abcxyz")
			}
			r := []ReceiverRecord{
				{
					PrivateKey:    key,
					TokenProvider: StaticTokenProvider("abcxyz"),
				},
			}
			server := Server{
				VoteHandler: vl,
				Records:     r,
			}
			go server.Serve(listener) //nolint:errcheck

			err = client.SendVote(v)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestServerv2Panic(t *testing.T) {
	expectedErr := errors.New("test error")
	vl := func(rv *Vote, ver Protocol) error {
		return expectedErr
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	defer listener.Close()
	r := []ReceiverRecord{
		{
			PrivateKey:    nil,
			TokenProvider: StaticTokenProvider("abcxyz"),
		},
	}
	server := Server{
		VoteHandler: vl,
		Records:     r,
	}
	go server.Serve(listener) //nolint:errcheck

	vote := Vote{
		ServiceName: "golang",
		Username:    "golang",
		Address:     "127.0.0.1",
	}
	client := NewV2Client(listener.Addr().String(), "abcxyz")
	err = client.SendVote(vote)
	if err == nil {
		t.Error("expected error, but didn't get any")
	}

	if errors.Is(err, expectedErr) {
		t.Errorf("expected error %q, but got %q", expectedErr, err)
	}
}

package votifier

import (
	"reflect"
	"testing"
	"time"
)

func TestEncodeV2(t *testing.T) {
	v := Vote{
		ServiceName: "golang",
		Username:    "golang",
		Address:     "127.0.0.1",
	}

	// Try to encrypt this vote.
	s, err := v.EncodeV2("abcxyz", "xyz")
	if err != nil {
		t.Error(err)
		return
	}

	// Try to decrypt this vote.
	var d Vote
	err = d.DecodeV2(s, StaticTokenProvider("abcxyz"), "xyz")
	if err != nil {
		t.Error(err)
		return
	}

	v.Timestamp = v.Timestamp.Round(time.Second)
	d.Timestamp = d.Timestamp.Round(time.Second)
	if !reflect.DeepEqual(v, d) {
		t.Error("votes don't match: ", v, "-", d)
	}
}

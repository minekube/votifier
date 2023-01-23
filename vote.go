package votifier

import (
	"time"
)

// Vote represents a Votifier vote.
type Vote struct {
	// The name of the service the user is voting from.
	ServiceName string `json:"serviceName"`

	// The user's Minecraft username.
	Username string `json:"username"`

	// The voting user's IP address.
	Address string `json:"address"`

	// The timestamp this vote was issued.
	Timestamp time.Time `json:"timeStamp"`
}

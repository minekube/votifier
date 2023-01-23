package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"log"
	"os"

	"go.minekube.com/votifier"
)

var (
	address     = flag.String("address", ":8192", "what host and port to connect to")
	keyFile     = flag.String("key", "", "key file to use")
	serviceName = flag.String("service", "go-votifier", "service name to use")
	username    = flag.String("user", "golang", "username to use")
	vAddress    = flag.String("user-address", "127.0.0.1", "address to use")
)

func main() {
	flag.Parse()

	file, err := os.ReadFile(*keyFile)
	if err != nil {
		log.Fatalf("loading public key: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(file))
	if err != nil {
		log.Fatalf("decoding public key: %v", err)
	}
	pkt, err := x509.ParsePKIXPublicKey(decoded)
	if err != nil {
		log.Fatalf("deserializing public key: %v", err)
	}

	key := pkt.(*rsa.PublicKey)
	client := votifier.NewV1Client(*address, key)
	v := votifier.Vote{
		ServiceName: *serviceName,
		Username:    *username,
		Address:     *vAddress,
	}
	err = client.SendVote(v)
	if err != nil {
		log.Fatalf("Failed to send vote: %v", err)
	}

	log.Println("Vote sent!")
}

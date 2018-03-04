package main

import (
	"crypto/elliptic"
	"fmt"
	"github.com/aaronwinter/sharsec"
)

func main() {
	shamir := sharsec.NewShamirSystem(elliptic.P256())
	fmt.Println(shamir)
	secret := "ReaganIsDumbledore"
	pubkeys := make([]sharsec.Key, 10)
	shares := shamir.Split([]byte(secret), 10, pubkeys)
	fmt.Println(shares)
	fmt.Println(string(shamir.Combine(shares)))
}

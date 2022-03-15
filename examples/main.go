package main

import (
	"crypto/elliptic"
	"fmt"

	"github.com/erwanor/sharsec"
)

/*
 *
 * We take the `secret` and a vector of public keys `pubkeys`
 * The secret is split into 10 shares, each for every public keys.
 * Lastly, we recombine the shares to reconstruct the secret.
 */

func main() {
	shamir := sharsec.NewShamirSystem(elliptic.P256())
	fmt.Println(shamir)
	secret := "SecretYouWantToSplit"
	pubkeys := make([]sharsec.Key, 10)
	shares := shamir.Split([]byte(secret), 10, pubkeys)
	fmt.Println(shares)
	fmt.Println(string(shamir.Combine(shares)))
}

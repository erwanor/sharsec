package sharsec

import (
	"math/big"

	"github.com/erwanor/sharsec/curvewrapper"
)

//   __      ___   ___ _  _ ___ _  _  ___
//   \ \    / /_\ | _ \ \| |_ _| \| |/ __|
//    \ \/\/ / _ \|   / .` || || .` | (_ |
//     \_/\_/_/ \_\_|_\_|\_|___|_|\_|\___|
//
// DO NOT USE THIS CODE IN PRODUCTION OR GET PWND
// NOT USING CONSTANT-TIME ARITHMETIC OPERATIONS
// NOT AUDITED

// THIS IS JUST TO MESS AROUND. NEVER LET THIS CODE HIT PRODUCTION

// SSS defines the behavior of a Secret Sharing Scheme (SSS)

type Key struct {
	priv []byte
	pub  curvewrapper.Point
}

type Share interface {
	Encrypt(Key) Share
	Decrypt(Key) Share
}

type SSS interface {
	Split([]byte, int, []Key) []Share
	Combine([]Share) []byte
}

type Polynomial interface {
	Eval(*big.Int, *big.Int) *big.Int
	String() string
}

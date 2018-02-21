package sharsec

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"sharsec/curvewrapper"
	"sharsec/finitefield"
)

type Shamir struct {
	ec elliptic.Curve
}

// For future iterations.
type EncryptedShare struct {
	SID   *finitefield.FpInt
	Value curvewrapper.Point
}

type ClearShare = EncryptedShare

func NewShamirSystem(curve elliptic.Curve) *Shamir {
	return &Shamir{
		ec: curve,
	}
}

func (s *Shamir) NewKey() Key {
	var k Key
	k.priv, k.pub.X, k.pub.Y, _ = elliptic.GenerateKey(s.ec, rand.Reader)
	k.pub.Curve = s.ec
	return k
}


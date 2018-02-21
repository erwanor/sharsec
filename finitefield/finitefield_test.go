package finitefield

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func assertBigInt(t *testing.T, a, b *big.Int, msgAndArgs ...interface{}) bool {
	return assert.Equal(t, a.Cmp(b), 0, msgAndArgs...)
}

func TestNewFpInt(t *testing.T) {
	// p = 2
	a := NewFpInt(big.NewInt(17), big.NewInt(2))
	assertBigInt(t, a.Value, big.NewInt(1), "17 reduce to 1 in F_2")

	b := NewFpInt(big.NewInt(256), big.NewInt(2))
	assertBigInt(t, b.Value, big.NewInt(0), "256 reduce to 0 in F_2")

	// p = 7
	d := NewFpInt(big.NewInt(8), big.NewInt(7))
	assertBigInt(t, d.Value, big.NewInt(1), "8 reduce to 1 in F_7")

	e := NewFpInt(big.NewInt(6), big.NewInt(7))
	assertBigInt(t, e.Value, big.NewInt(6), "6 reduce to 6 in F_7")

	f := NewFpInt(big.NewInt(14), big.NewInt(7))
	assertBigInt(t, f.Value, big.NewInt(0), "14 reduce to 0 in F_7")

	g := NewFpInt(big.NewInt(25), big.NewInt(7))
	assertBigInt(t, g.Value, big.NewInt(4), "25 reduce to 4 in F_7")
}

func TestFpInt_Add(t *testing.T) {
	// p = 2
	a := NewFpInt(big.NewInt(1), big.NewInt(2))
	b := NewFpInt(big.NewInt(2), big.NewInt(2))
	c := NewFpInt(big.NewInt(0), big.NewInt(2))
	c.Add(a, b)
	//if err != nil {
	//	t.Fatal(err)
	//}
	assertBigInt(t, c.Value, big.NewInt(1), "1 add 2 reduce to 1 in F_2")

	d := NewFpInt(big.NewInt(5189), big.NewInt(2))
	c.Add(d, b)
	//if err != nil {
	//	t.Fatal(err)
	//}
	assertBigInt(t, c.Value, big.NewInt(1), "5189 add 2 reduce to 1 in F_2")

	f := NewFpInt(big.NewInt(147), big.NewInt(2))
	c.Add(f, a)
	//if err != nil {
	//	t.Fatal(err)
	//}
	assertBigInt(t, c.Value, big.NewInt(0), "147 add 1 reduce to 0 in F_2")
}

func TestFpInt_Mul(t *testing.T) {
	// p = 2
	a := NewFpInt(big.NewInt(1), big.NewInt(2))
	b := NewFpInt(big.NewInt(0), big.NewInt(2))
	c := NewFpInt(big.NewInt(0), big.NewInt(2))
	c.Mul(a, b)
	//if err != nil {
	//	t.Fatal(err)
	//}
	assertBigInt(t, c.Value, big.NewInt(0), "1 mul 0 reduce to 0 in F_2")

	c.Mul(a, a)
	//if err != nil {
	//	t.Fatal(err)
	//}
	assertBigInt(t, c.Value, big.NewInt(1), "1 mul 1 reduce to 1 in F_2")

	c.Mul(b, b)
	//if err != nil {
	//	t.Fatal(err)
	//}
	assertBigInt(t, c.Value, big.NewInt(0), "0 mul 0 reduce to 0 in F_2")
}

package effects

import (
	"crypto/rand"
	"math/big"
)

// Returns a random number between 0.0 and 1.0.
func randomNonNegative() float64 {
	const long = 4_294_967_296
	n, err := rand.Int(rand.Reader, big.NewInt(long))
	if err != nil {
		return 0.0
	}
	return float64(n.Int64()) / long
}

// Returns a random number between -1.0 and 1.0.
func random() float64 {
	const slong = 2_147_483_648
	n, err := rand.Int(rand.Reader, big.NewInt(slong))
	if err != nil {
		return 0.0
	}
	return float64(n.Int64())/slong - 1
}

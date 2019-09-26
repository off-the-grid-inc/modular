package modular

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type Int struct {
	Value *big.Int
	Base  *big.Int
}

// Default prime (256-bit secp256k1 EC order)
var (
	defaultP, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	Zero        = big.NewInt(0)
)

// ARITHMETIC OPERATIONS

// Add - Modular addition of (an arbitrary number of) field elements
func (n *Int) Add(x, y *Int) *Int {
	v := new(big.Int).Add(x.Value, y.Value)
	n.Value = v.Mod(v, x.Base)
	n.Base = x.Base
	return n
}

// Mul - Modular multiplication of two field elements
func (n *Int) Mul(x, y *Int) *Int {
	v := new(big.Int).Mul(x.Value, y.Value)
	n.Value = v.Mod(v, x.Base)
	n.Base = x.Base
	return n
}

// Sub - Modular Subtraction of two field elements
func (n *Int) Sub(x, y *Int) *Int {
	v := new(big.Int).Sub(x.Value, y.Value)
	n.Value = v.Mod(v, x.Base)
	n.Base = x.Base
	return n
}

// LinearCombination is the dot product of two vectors in a finite field. Vectors should have the same dimensionality for proper behavior
func (n *Int) LinearCombination(vec1, vec2 []*Int) *Int {
	v := big.NewInt(0)
	n.Base = vec1[0].Base
	for i := range vec1 {
		v.Add(v, new(big.Int).Mul(vec1[i].Value, vec2[i].Value))
	}
	n.Value = v.Mod(v, n.Base)
	return n
}

// ModInverse is a custom implementation for the modular inverse of a field element
func ModInverse(number *Int) *Int {
	copy := big.NewInt(0).Set(number.Value)
	pcopy := big.NewInt(0).Set(number.Base)
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set(number.Base)

	result.Add(result, y)
	result.Mod(result, number.Base)
	return &Int{
		Value: result,
		Base:  number.Base,
	}
}

// Exp - exponentiate in a finite field
func (n *Int) Exp(base, exp *Int) *Int {
	n.Value = new(big.Int).Exp(base.Value, exp.Value, base.Base)
	n.Base = base.Base
	return n
}

// Cmp compares two finite field elements
func (n *Int) Cmp(x *Int) int {
	return n.Value.Cmp(x.Value)
}

// Bytes returs the byte array representation of a field element
func (n *Int) Bytes() []byte {
	return n.Value.Bytes()
}

// String returns the string representation
func (n *Int) String() string {
	return n.Value.String()
}

// Helpers

// Note: when instansiating new Ints the prime is assigned as a pointer which is assumed to be static.
// If Int.Base prime pointers are altered it could cause strange behaviour and/or data races.

// NewInt creates a modular Int from int64
func NewInt(i int64, prime *big.Int) *Int {
	num := big.NewInt(i)
	if prime == nil {
		prime = defaultP
	}
	return &Int{
		Value: num.Mod(num, prime),
		Base:  prime,
	}
}

// SetFromBytes creates a modular Int from a byte array
func IntFromBytes(b []byte, prime *big.Int) *Int {
	num := new(big.Int).SetBytes(b)
	if prime == nil {
		prime = defaultP
	}
	return &Int{
		Value: num.Mod(num, prime),
		Base:  prime,
	}
}

// SetFromBig creates a modular Int from a big Int
func IntFromBig(num *big.Int, prime *big.Int) *Int {
	if prime == nil {
		prime = defaultP
	}
	return &Int{
		Value: new(big.Int).Mod(num, prime),
		Base:  prime,
	}
}

// SetFromString creates a modular Int from a string representation of an integer in a specific Base.
func IntFromString(str string, Base int, prime *big.Int) (*Int, error) {
	num, ok := new(big.Int).SetString(str, Base)
	if !ok {
		return nil, errors.New("Cannot interpret string")
	}
	if prime == nil {
		prime = defaultP
	}
	return &Int{
		Value: new(big.Int).Mod(num, prime),
		Base:  prime,
	}, nil
}

// RandInt creates a new random modular Int within the range [0, p)
func RandInt(prime *big.Int) (*Int, error) {
	if prime == nil {
		prime = defaultP
	}
	max := new(big.Int).Sub(prime, big.NewInt(1))
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return &Int{
		Value: result,
		Base:  prime,
	}, nil
}

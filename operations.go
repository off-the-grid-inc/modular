package modular

import (
	"math/big"
	"crypto/rand"
	"errors"
)

type Int big.Int

// Default prime (256-bit secp256k1 EC order)
var (
	p, _ = IntFromString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
)

// SetP resets the (prime) modulus across the package
func SetP(new_prime *Int) *Int {
	p = new_prime
	return p
}

func GetP() *Int {
	return p
}

// ARITHMETIC OPERATIONS

// Add - Modular addition of (an arbitrary number of) field elements
func (n *Int) Add(nums ...*Int) *Int{
	out := big.NewInt(0)
	for _, n := range nums {
		out.Add(out, (*big.Int)(n))
	} 
	out.Mod(out, (*big.Int)(p))
	*n = (Int)(*out)
	return n
}

// Mul - Modular multiplication of two field elements
func (n *Int) Mul(x, y *Int) *Int {
	out := new(big.Int).Mul((*big.Int)(x), (*big.Int)(y))
	out.Mod(out, (*big.Int)(p))
	*n = (Int)(*out)
	return n
}

// Sub - Modular Subtraction of two field elements
func (n *Int) Sub(x, y *Int) *Int {
	out := new(big.Int).Sub((*big.Int)(x), (*big.Int)(y))
	out.Mod(out, (*big.Int)(p))
	*n = (Int)(*out)
	return n
}

// LinearCombination is the dot product of two vectors in a finite field. Vectors should have the same dimensionality for proper behavior
func (n *Int) LinearCombination(vec1 []*Int, vec2 []*Int) *Int {
	out := big.NewInt(0)
	for i := range vec1 {
		out.Add(out, new(big.Int).Mul((*big.Int)(vec1[i]), (*big.Int)(vec2[i])))
	}
	out.Mod(out, (*big.Int)(p))
	*n = (Int)(*out)
	return n
}

// ModInverse is a custom implementation for the modular inverse of a field element
func ModInverse(number *Int) *Int {
	copy := big.NewInt(0).Set((*big.Int)(number))
	pcopy := big.NewInt(0).Set((*big.Int)(p))
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set((*big.Int)(p))

	result.Add(result, y)
	result.Mod(result, (*big.Int)(p))
	return (*Int)(result)
}

// Exp - exponentiate in a finite field
func (n *Int) Exp(base, exp *Int) *Int {
	*n = (Int)(*new(big.Int).Exp((*big.Int)(base), (*big.Int)(exp), (*big.Int)(p)))
	return n
}

// Cmp compares two finite field elements
func (n *Int) Cmp(x *Int) int {
	return (*big.Int)(n).Cmp((*big.Int)(x))
}

// Bytes returs the byte array representation of a field element
func (n *Int) Bytes() []byte {
	return (*big.Int)(n).Bytes()
}


// HELPERS

// Note: New Int's are not automatically reduced mod P

// NewInt creates a modular Int from int64 
func NewInt(i int64) *Int {
	num := big.NewInt(i)
	return (*Int)(num)
}

// SetFromBytes creates a modular Int from a byte array
func IntFromBytes(b []byte) *Int {
	num := new(big.Int).SetBytes(b)
	return (*Int)(num)
}

// SetFromBig creates a modular Int from a big Int
func IntFromBig(num *big.Int) *Int {
	return (*Int)(new(big.Int).Set(num))
}

// SetFromString creates a modular Int from a string representation of an integer in a specific base.
func IntFromString(str string, base int) (*Int, error) {
	num, err := new(big.Int).SetString(str, base)
	if !err {
		return nil, errors.New("Could not set string")
	}
	return (*Int)(num), nil
}

// IsModP checks whether a modular Int actually lies within the field order
func (n *Int) IsModP() bool {
	if n.Cmp(p) == -1 {
		return true
	}
	return false
}

func (n *Int) AsBig() *big.Int {
	return (*big.Int)(n)
}

func (n *Int) Mod() *Int {
	*n = (Int)(*new(big.Int).Mod((*big.Int)(n), (*big.Int)(p)))
	return n
}

// RandInt creates a new random modular Int within the range [0, p)
func RandInt() (*Int, error) {
	max := big.NewInt(0).Set((*big.Int)(p))
	max.Sub(max, big.NewInt(1))
	result, err := rand.Int(rand.Reader, max)
	if err != nil{
		return nil, err
	}
	return (*Int)(result), nil
}

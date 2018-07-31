package modular

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModInverse(t *testing.T){
	require := require.New(t)

	for i := 0; i < 100; i++ {
		point, err := RandInt()
		require.NoError(err)
		require.Equal(point.Cmp(p), -1, "out of bounds")
		inverse := ModInverse(point)
		require.Equal(inverse.Cmp(p), -1, "Inverse out of bounds")
		point.Mul(point, inverse)		
		require.Equal(0, point.Cmp(NewInt(1)),  "Inverse incorrect")
	}
}
func TestRandomInt(t *testing.T) {
	require := require.New(t)	
	for i := 0; i < 100; i++ {
		x, err := RandInt()
		require.NoError(err)
		require.Equal(-1, x.Cmp(p), "Inverse out of bounds")
	}
}

func TestOperations(t *testing.T) {
	require := require.New(t)

	// basic multiplication
	check := big.NewInt(1234)
	check.Mul(check, big.NewInt(2))
	res := new(Int).Mul(NewInt(1234), NewInt(2))
	require.Equal(0, res.Cmp((*Int)(check)), "multiplication failure")

	// modular overflow multiplication
	check = new(big.Int).Exp(big.NewInt(2), big.NewInt(500), nil)
	res = new(Int).Exp(NewInt(2), NewInt(500))
	check.Mul(check, check)
	res.Mul(res, res)
	require.Equal(-1, res.Cmp((*Int)(check)), "did not automatically reduce")
	check.Mod(check, (*big.Int)(p))
	require.Equal(0, res.Cmp((*Int)(check)), "did not reduce properly")

	// basic addition
	check = big.NewInt(1234)
	check.Mul(check, big.NewInt(3))
	res = new(Int).Add(NewInt(1234), NewInt(1234), NewInt(1234))
	require.Equal(0, res.Cmp((*Int)(check)), "addition failure")

	// modular overflow addition
	check = new(big.Int).Exp(big.NewInt(2), big.NewInt(1000), nil)
	check.Mul(check, big.NewInt(3))
	res = new(Int).Exp(NewInt(2), NewInt(1000))
	res.Add(res, res, res)
	require.Equal(-1, res.Cmp((*Int)(check)), "did not automatically reduce")
	check.Mod(check, (*big.Int)(p))	
	require.Equal(0, res.Cmp((*Int)(check)), "did not reduce properly")

	// modular subtraction
	check = new(big.Int).Exp(big.NewInt(2), big.NewInt(1000), nil)
	check.Sub(check, big.NewInt(1000000))
	res = new(Int).Exp(NewInt(2), NewInt(1000))
	res.Sub(res, NewInt(1000000))
	require.Equal(-1, res.Cmp((*Int)(check)), "did not automatically reduce")
	check.Mod(check, (*big.Int)(p))	
	require.Equal(0, res.Cmp((*Int)(check)), "did not reduce properly")

	// test linear combination
	vec1 := []*Int{NewInt(1), NewInt(2), NewInt(3), NewInt(4)}
	vec2 := []*Int{NewInt(18000), NewInt(9000), NewInt(6000), NewInt(4500)}
	check = new(big.Int).Mul(big.NewInt(18000), big.NewInt(4))
	res = new(Int).LinearCombination(vec1, vec2)
	require.Equal(0, res.Cmp((*Int)(check)), "linear combination failure")
}

func TestChangePrime(t *testing.T) {
	require := require.New(t)
	newp, err := IntFromString("9049555721791567387589049441228905843", 10)
	require.NoError(err)

	SetP(newp)
	require.Equal(0, p.Cmp(newp), "change global prime failed")
}


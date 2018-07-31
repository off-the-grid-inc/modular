package modular

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicMatrix(t *testing.T){
	require := require.New(t)

	data := make([]*Int, 10)
	for i := range data {
		v, err := RandInt()
		require.NoError(err)
		data[i] = v
	}

	// Test Create Matrix
	m := NewMatrix(2, 5, data)

	// Test Get Row/Col
	row := m.GetRow(2)
	require.Equal(row[4].Cmp(data[9]), 0, "get row failed")
	col := m.GetCol(5)
	require.Equal(row[4].Cmp(col[1]), 0, "get column failed")

	// Test Scalar Multiplication
	m.ScalarMul(new(Int).Exp(NewInt(2), NewInt(256)))
	require.Equal(0, m.values[9].Cmp(data[9].Mul(data[9], new(Int).Exp(NewInt(2), NewInt(256)))), "scalar mult failed")

	// Test Set Row/Col
	m.SetRow(1, []*Int{NewInt(1),NewInt(1),NewInt(1),NewInt(1),NewInt(1)})
	require.Equal(0, m.values[0].Cmp(m.values[4]), "set row failed")
	m.SetCol(5, []*Int{NewInt(1), NewInt(1)})
	require.Equal(0, m.values[9].Cmp(m.values[4]), "set column failed")

}

func TestMultiplication(t *testing.T) {
	require := require.New(t)

	m1 := NewMatrix(2, 3, []*Int{NewInt(1), NewInt(2), NewInt(3), NewInt(4), NewInt(5), NewInt(6)})
	m2 := NewMatrix(3, 1, []*Int{NewInt(3), NewInt(2), NewInt(1)})
	res, err := new(Matrix).Mul(m1, m2)
	require.NoError(err)
	require.Equal(len(res.values), 2, "wrong structure")
	require.Equal(0, res.values[0].Cmp(NewInt(10)), "multiplication failed")
	require.Equal(0, res.values[1].Cmp(NewInt(28)), "multiplication failed")	
}

func TestInverse(t *testing.T) {
	require := require.New(t)

	// Test Gauss Jordan
	linearSystem := NewMatrix(2, 2, []*Int{NewInt(1), NewInt(3), NewInt(2), NewInt(7)})
	linearSystemResult := []*Int{NewInt(1), NewInt(0)}
	result, err := GaussJordan(linearSystem.Represent2D(), linearSystemResult)
	require.NoError(err, "gauss jordan failed")
	require.Equal(0, result[0].Cmp(NewInt(7)), "gauss jordan failed")
	require.Equal(0, result[1].Cmp(NewInt(-2).Mod()), "gauss jordan failed")

	// Test Inverses
	m := NewMatrix(2, 2, []*Int{NewInt(7), NewInt(-3).Mod(), NewInt(-2).Mod(), NewInt(1)})
	inv, err := m.Inverse()
	require.NoError(err, "inverse failed")
	require.Equal(0, inv.values[0].Cmp(NewInt(1)), "inverse failed")
	
}


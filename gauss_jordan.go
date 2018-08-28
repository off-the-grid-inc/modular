package modular

// repurposed from github.com/thunpin/go-math/

import (
	"errors"
	"fmt"
	"math/big"
)

func PrintArray(arr []*Int) {
	for _, a := range arr {
		b := (*big.Int)(a)
		fmt.Printf("%d, ", b)
	}
	fmt.Println("")
}

func Max(arr []*Int) (int, *Int, error) {
	if len(arr) == 0 {
		return 0, nil, errors.New("Empty array is not valid")
	}
	max := arr[0]
	idx := int(0)
	for i, a := range arr {
		if a.Cmp(max) == 1 {
			max = a
			idx = i
		}
	}
	return idx, max, nil
}

// Appends linearSystemResult as the last column to linearSystem
func createExtendedArray(linearSystem [][]*Int, linearSystemResult []*Int, n int) [][]*Int {
	nrows := len(linearSystem)
	ncols := len(linearSystem[0])

	extended := make([][]*Int, nrows)
	for i := 0; i < nrows; i++ {
		extended[i] = make([]*Int, ncols+1)

		for j := 0; j < ncols; j++ {
			extended[i][j] = linearSystem[i][j]
		}
		extended[i][ncols] = linearSystemResult[i]
	}

	return extended
}

func ExtractColumn(matrix [][]*Int, j int) []*Int {
	column := make([]*Int, len(matrix))
	for i, a := range matrix {
		column[i] = a[j]
	}
	return column
}

func SwapRows(matrix [][]*Int, i int, j int) {
	// tmp := matrix[i]
	// matrix[i] = matrix[j]
	// matrix[j] = tmp
	matrix[i], matrix[j] = matrix[j], matrix[i]
}

// GaussJordan solve a linear system using Gauss Jordan algorithm
func GaussJordan(linearSystem [][]*Int, linearSystemResult []*Int) ([]*Int, error) {
	// TODO: Validate system
	// TODO: Handle no solutions
	nrows := len(linearSystem)
	ncols := len(linearSystem[0])
	extendedMatrix := createExtendedArray(linearSystem, linearSystemResult, ncols)

	ncols++ // Matrix has been extended
	//	fmt.Printf("sizes: %d, %d\n", ncols, len(extendedMatrix[0]))
	h := 0
	k := 0
	column := make([]*Int, len(extendedMatrix))
	for (h < nrows) && (k < ncols) {
		/* Find the k-th pivot: */
		//		fmt.Printf("%d, %d\n", h, k)
		column = ExtractColumn(extendedMatrix, k)
		i_max, max, err := Max(column[h:])
		if err != nil {
			return nil, err
		}
		// Index is relative, we make it absolute
		i_max = i_max + h
		if max.Cmp(NewInt(0)) == 0 {
			/* No pivot in this column, pass to next column */
			// fmt.Println("DONT ENTER HERE")
			k++
		} else {
			SwapRows(extendedMatrix, h, i_max)
			/* Do for all rows below pivot: */
			for i := h + 1; i < nrows; i++ {
				factor := new(Int).Mul(extendedMatrix[i][k], ModInverse(extendedMatrix[h][k]))
				/* Fill with zeros the lower part of pivot column: */
				extendedMatrix[i][k] = NewInt(0)
				/* Do for all remaining elements in current row: */
				for j := k + 1; j < ncols; j++ {
					num := new(Int).Mul(extendedMatrix[h][j], factor)
					extendedMatrix[i][j] = new(Int).Sub(extendedMatrix[i][j], num)
				}
			}
			/* Increase pivot row and column */
			h++
			k++
		}
	}

	// The matrix is upper triangular, but not in row echelon form yet

	fmt.Println("")
	fmt.Println("EXTENDED MATRIX")
	for _, a := range extendedMatrix {
		PrintArray(a)
	}
	fmt.Println("")

	pivots := make([]int, nrows)
	for i := 0; i < nrows; i++ {
		pivots[i] = -1
	}
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			if extendedMatrix[i][j].Cmp(NewInt(0)) != 0 {
				pivots[i] = j
				break
			}
		}
		if pivots[i] == ncols-1 {
			errors.New("There are no solutions")
		}
	}
	fmt.Println("Pivots: ")
	for _, p := range pivots {
		fmt.Printf("%d ", p)
	}
	fmt.Println("")

	result := make([]*Int, ncols-1)
	for j := 0; j < ncols-1; j++ {
		result[j] = NewInt(0)
	}
	for i := nrows - 1; i >= 0; i-- {
		if pivots[i] != -1 {
			// fmt.Printf("DAMN i: %d\n", i)
			result[pivots[i]] = extendedMatrix[i][ncols-1]
			for i_ := nrows - 1; i_ > i && pivots[i_] != -1; i_-- {
				mult := new(Int).Mul(result[pivots[i_]], extendedMatrix[i][pivots[i_]])
				result[pivots[i]] = new(Int).Sub(result[pivots[i]], mult)
				result[pivots[i]] = new(Int).Mul(result[pivots[i]], ModInverse(extendedMatrix[i][pivots[i]]))
			}
		}
	}
	return result, nil
}

func GaussJordan2(linearSystem [][]*Int, linearSystemResult []*Int) ([]*Int, error) {

	result := make([]*Int, len(linearSystem))
	// Check sizes and non-zero columns
	err := validateLinearSystemMatrix(linearSystem, linearSystemResult)
	if err != nil {
		return result, err
	}

	n := len(linearSystem)
	// Append result vector to last column
	extendsdMatrix := createExtendedArray(linearSystem, linearSystemResult, n)
	// try put 1 in all a[i,j] where i == j
	for i := 0; i < n; i++ {
		err = solveLineGaussJordan(extendsdMatrix, i, n, n+1)
		if err != nil {
			break
		}
	}

	if err == nil {
		for i := 0; i < n; i++ {
			result[i] = extendsdMatrix[i][n]
		}
	}

	return result, err
}

func solveLineGaussJordan(linearSystem [][]*Int, i int, n int, m int) error {
	var err error

	// if the current diagonal is 0 we need swap then
	if linearSystem[i][i].Cmp(NewInt(0)) == 0 {
		if i == (n - 1) {
			err = ErrDontExistSolution
		} else {
			err = swapGaussJordan(linearSystem, i, n)
		}
	}

	if err == nil {
		factor := ModInverse(linearSystem[i][i])
		for j := 0; j < m; j++ {
			if j == i {
				linearSystem[i][j] = NewInt(1)
			} else {
				linearSystem[i][j].Mul(linearSystem[i][j], factor)
			}
		}

		// put 0 in all a[j,i]
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}

			factor := new(Int).Mul(linearSystem[j][i], ModInverse(linearSystem[i][i]))
			for k := 0; k < m; k++ {
				if k == i {
					linearSystem[j][k] = NewInt(0)
				} else {
					num := new(Int).Mul(linearSystem[i][k], factor)
					linearSystem[j][k].Sub(linearSystem[j][k], num)
				}
			}
		}
	}

	return err
}

// Looks for a row j below the i-th row such that the entry [j][i] is non-zero.
// This is only called if [i][i] was zero to begin with
// TODO: I think this is the problem! A system with multiple solutions will have zero rows at the bottom so no swap operation will be found
func swapGaussJordan(linearSystem [][]*Int, i int, n int) error {
	var err error

	swaped := false
	for j := i + 1; j < n; j++ {
		if linearSystem[j][i].Cmp(NewInt(0)) != 0 {
			swaped = true
			tmp := linearSystem[i]
			linearSystem[i] = linearSystem[j]
			linearSystem[j] = tmp

			break
		}
	}

	if swaped == false {
		err = ErrDontExistSolution
	}

	return err
}

// HELPERS

var (
	ErrInvalidLinearSystemMatrix       = errors.New("the linear system matrix is not a quadratic matrix")
	ErrInvalidLinearSystemResultMatrix = errors.New("the linear system result matrix is invalid")
	ErrDontExistSolution               = errors.New("dont exist solution to this linear system")
)

// validateLinearSystemMatrix verify if the linearSystem is a quadratic matrix n
// and the linearSystemResult is matrix n x 1
func validateLinearSystemMatrix(linearSystem [][]*Int, linearSystemResult []*Int) error {

	var err error
	n := len(linearSystem)
	if n == 0 {
		err = ErrInvalidLinearSystemMatrix
	}

	l := len(linearSystem[0])

	// verify if the linear system matrix is quadratic
	// TODO: We must allow non-square systems (btw, square, not 'quadratic'!)
	for i := 0; i < n; i++ {
		m := len(linearSystem[i])
		if m != l {
			err = ErrInvalidLinearSystemMatrix
			break
		}
	}

	// verify if the linear system matrix and your result have the same number of rows
	if err == nil && len(linearSystemResult) != n {
		err = ErrInvalidLinearSystemResultMatrix
	}

	if err == nil && canFoundSolution(linearSystem, n) == false {
		err = ErrDontExistSolution
	}

	return err
}

// TODO: This is checking that the matrix does not have a zero column. This may be a problem
// canFoundSolution verify if is possible found a solution
func canFoundSolution(linearSystem [][]*Int, dimension int) bool {
	result := true

	// verify if exist a "i" column different of zero
	for j := 0; j < dimension; j++ {
		result = false
		for i := 0; i < dimension; i++ {
			if linearSystem[i][j].Cmp(NewInt(0)) != 0 {
				result = true
				continue
			}
		}

		if result == false {
			break
		}
	}

	return result
}

package modular

import (
	"errors"
	// "fmt"
	// "math/big"
)

func NonZeroEntry(arr []*Int) (int, *Int, error) {
	if len(arr) == 0 {
		return 0, nil, errors.New("Empty array is not valid")
	}
	nz := arr[0]
	idx := int(0)
	for i, a := range arr {
		if a.Cmp(NewInt(0)) != 0 {
			nz = a
			idx = i
			break
		}
	}
	return idx, nz, nil
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

// Extracts column j from the matrix
func ExtractColumn(matrix [][]*Int, j int) []*Int {
	column := make([]*Int, len(matrix))
	for i, a := range matrix {
		column[i] = a[j]
	}
	return column
}

// Swaps rows i and j of matrix
func SwapRows(matrix [][]*Int, i int, j int) {
	matrix[i], matrix[j] = matrix[j], matrix[i]
}

// GaussJordan solve a linear system using Gauss Jordan algorithm
func GaussJordan(linearSystem [][]*Int, linearSystemResult []*Int) ([]*Int, error) {
	nrows := len(linearSystem)
	ncols := len(linearSystem[0])
	dresult := len(linearSystemResult)
	if (nrows*ncols == 0) || (ncols != dresult) {
		return nil, errors.New("Wrong dimensions for the system")
	}

	extendedMatrix := createExtendedArray(linearSystem, linearSystemResult, ncols)

	ncols++ // Matrix has been extended

	h := 0
	k := 0
	column := make([]*Int, len(extendedMatrix))

	for (h < nrows) && (k < ncols) {
		/* Find the k-th pivot: */
		column = ExtractColumn(extendedMatrix, k)

		i_max, nz, err := NonZeroEntry(column[h:])

		if err != nil {
			return nil, err
		}
		// Index is relative, we make it absolute
		i_max = i_max + h
		if nz.Cmp(NewInt(0)) == 0 {
			/* No pivot in this column, pass to next column */
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
	// List holding the pivots positions for each row
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
			err := errors.New("There are no solutions")
			return nil, err
		}
	}

	result := make([]*Int, ncols-1)
	for j := 0; j < ncols-1; j++ {
		result[j] = NewInt(0)
	}

	// Calculate result in position pivots[i]
	for i := nrows - 1; i >= 0; i-- {
		if pivots[i] != -1 {
			result[pivots[i]] = extendedMatrix[i][ncols-1]
			for i_ := nrows - 1; i_ > i; i_-- {
				if pivots[i_] != -1 {
					mult := new(Int).Mul(result[pivots[i_]], extendedMatrix[i][pivots[i_]])

					result[pivots[i]] = new(Int).Sub(result[pivots[i]], mult)
				}

			}

			result[pivots[i]] = new(Int).Mul(result[pivots[i]], ModInverse(extendedMatrix[i][pivots[i]]))
		}
	}
	return result, nil
}

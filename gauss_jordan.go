package modular

// repurposed from github.com/thunpin/go-math/

import "errors"

// GaussJordan solve a linear system using Gauss Jordan algorithm
func GaussJordan(linearSystem [][]*Int, linearSystemResult []*Int) ([]*Int, error) {

	result := make([]*Int, len(linearSystem))
	err := validateLinearSystemMatrix(linearSystem, linearSystemResult)
	if err != nil {
		return result, err
	}

	n := len(linearSystem)
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
	ErrInvalidLinearSystemMatrix = errors.New("the linear system matrix is not a quadratic matrix")
	ErrInvalidLinearSystemResultMatrix = errors.New("the linear system result matrix is invalid")
	ErrDontExistSolution = errors.New("dont exist solution to this linear system")
)

// validateLinearSystemMatrix verify if the linearSystem is a quadratic matrix n
// and the linearSystemResult is matrix n x 1
func validateLinearSystemMatrix(linearSystem [][]*Int, linearSystemResult []*Int) error {

	var err error
	n := len(linearSystem)

	// verify if the linear system matrix is quadratic
	for i := 0; i < n; i++ {
		m := len(linearSystem[i])
		if m != n {
			err = ErrInvalidLinearSystemMatrix
			break
		}
	}

	// verify if the linear system matrix and your result has the same number of
	// lines
	if err == nil && len(linearSystemResult) != n {
		err = ErrInvalidLinearSystemResultMatrix
	}

	if err == nil && canFoundSolution(linearSystem, n) == false {
		err = ErrDontExistSolution
	}

	return err
}

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

func createExtendedArray(linearSystem [][]*Int, linearSystemResult []*Int, n int) [][]*Int {

	extended := make([][]*Int, n)
	for i := 0; i < n; i++ {
		extended[i] = make([]*Int, n+1)

		for j := 0; j < n; j++ {
			extended[i][j] = linearSystem[i][j]
		}
		extended[i][n] = linearSystemResult[i]
	}

	return extended
}
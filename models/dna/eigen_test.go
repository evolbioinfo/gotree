package dna

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

// Tests models with parameters making them equal to JC

func TestJCEigens(t *testing.T) {
	var v, l, r []float64
	var vMat, lMat, rMat mat.Matrix
	var err error
	var expQMatrix *mat.Dense
	var resQMatrix *mat.Dense = mat.NewDense(4, 4, nil)
	var sub *mat.Dense = mat.NewDense(4, 4, nil)

	expQMatrix = mat.NewDense(4, 4, []float64{
		-1.0000000, 0.3333333, 0.3333333, 0.3333333,
		0.3333333, -1.0000000, 0.3333333, 0.3333333,
		0.3333333, 0.3333333, -1.0000000, 0.3333333,
		0.3333333, 0.3333333, 0.3333333, -1.0000000,
	})

	m := NewJCModel()
	m.InitModel()
	if v, l, r, err = m.Eigens(); err != nil {
		t.Errorf("Error while computing JC eigen vectors: %v", err)
	}

	// Transpose because it was in col-major formrat
	lMat = mat.NewDense(4, 4, l).T()
	rMat = mat.NewDense(4, 4, r).T()

	vMat = mat.NewDiagDense(4, v)

	resQMatrix.Mul(lMat, vMat)
	resQMatrix.Mul(resQMatrix, rMat)
	sub.Sub(resQMatrix, expQMatrix)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if sub.At(i, j) > .0000001 {
				t.Errorf("Expected QMatrix is different from Resulting QMatrix %v", sub.At(i, j))
			}
		}
	}
}

func TestF81Eigens(t *testing.T) {
	var v, l, r []float64
	var vMat, lMat, rMat mat.Matrix
	var err error
	var expQMatrix *mat.Dense
	var resQMatrix *mat.Dense = mat.NewDense(4, 4, nil)
	var sub *mat.Dense = mat.NewDense(4, 4, nil)

	expQMatrix = mat.NewDense(4, 4, []float64{
		-1.0000000, 0.3333333, 0.3333333, 0.3333333,
		0.3333333, -1.0000000, 0.3333333, 0.3333333,
		0.3333333, 0.3333333, -1.0000000, 0.3333333,
		0.3333333, 0.3333333, 0.3333333, -1.0000000,
	})
	t.Logf("expQ = %v", mat.Formatted(expQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	m := NewF81Model()
	m.InitModel(1./4., 1./4., 1./4., 1./4.)

	if v, l, r, err = m.Eigens(); err != nil {
		t.Errorf("Error while computing F81 eigen vectors: %v", err)
	}

	// Transpose because it was in col-major formrat
	lMat = mat.NewDense(4, 4, l).T()
	rMat = mat.NewDense(4, 4, r).T()

	vMat = mat.NewDiagDense(4, v)

	resQMatrix.Mul(lMat, vMat)
	resQMatrix.Mul(resQMatrix, rMat)
	sub.Sub(resQMatrix, expQMatrix)

	t.Logf("resQ = %v", mat.Formatted(resQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if sub.At(i, j) > .0000001 {
				t.Errorf("Expected QMatrix is different from Resulting QMatrix %v", sub.At(i, j))
			}
		}
	}
}

func TestF84Eigens(t *testing.T) {
	var v, l, r []float64
	var vMat, lMat, rMat mat.Matrix
	var err error
	var expQMatrix *mat.Dense
	var resQMatrix *mat.Dense = mat.NewDense(4, 4, nil)
	var sub *mat.Dense = mat.NewDense(4, 4, nil)

	expQMatrix = mat.NewDense(4, 4, []float64{
		-1.0000000, 0.3333333, 0.3333333, 0.3333333,
		0.3333333, -1.0000000, 0.3333333, 0.3333333,
		0.3333333, 0.3333333, -1.0000000, 0.3333333,
		0.3333333, 0.3333333, 0.3333333, -1.0000000,
	})
	t.Logf("expQ = %v", mat.Formatted(expQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	m := NewF84Model()
	m.InitModel(0., 1./4., 1./4., 1./4., 1./4.)

	if v, l, r, err = m.Eigens(); err != nil {
		t.Errorf("Error while computing F81 eigen vectors: %v", err)
	}

	// Transpose because it was in col-major formrat
	lMat = mat.NewDense(4, 4, l).T()
	rMat = mat.NewDense(4, 4, r).T()

	vMat = mat.NewDiagDense(4, v)

	resQMatrix.Mul(lMat, vMat)
	resQMatrix.Mul(resQMatrix, rMat)
	sub.Sub(resQMatrix, expQMatrix)

	t.Logf("resQ = %v", mat.Formatted(resQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if sub.At(i, j) > .0000001 {
				t.Errorf("Expected QMatrix is different from Resulting QMatrix %v", sub.At(i, j))
			}
		}
	}
}

func TestGTREigens(t *testing.T) {
	var v, l, r []float64
	var vMat, lMat, rMat mat.Matrix
	var err error
	var expQMatrix *mat.Dense
	var resQMatrix *mat.Dense = mat.NewDense(4, 4, nil)
	var sub *mat.Dense = mat.NewDense(4, 4, nil)

	expQMatrix = mat.NewDense(4, 4, []float64{
		-1.0000000, 0.3333333, 0.3333333, 0.3333333,
		0.3333333, -1.0000000, 0.3333333, 0.3333333,
		0.3333333, 0.3333333, -1.0000000, 0.3333333,
		0.3333333, 0.3333333, 0.3333333, -1.0000000,
	})
	t.Logf("expQ = %v", mat.Formatted(expQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	m := NewGTRModel()
	m.InitModel(1., 1., 1., 1., 1., 1., 1./4., 1./4., 1./4., 1./4.)

	if v, l, r, err = m.Eigens(); err != nil {
		t.Errorf("Error while computing F81 eigen vectors: %v", err)
	}

	// Transpose because it was in col-major formrat
	lMat = mat.NewDense(4, 4, l).T()
	rMat = mat.NewDense(4, 4, r).T()

	vMat = mat.NewDiagDense(4, v)

	resQMatrix.Mul(lMat, vMat)
	resQMatrix.Mul(resQMatrix, rMat)
	sub.Sub(resQMatrix, expQMatrix)
	t.Logf("resQ = %v", mat.Formatted(resQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if sub.At(i, j) > .0000001 {
				t.Errorf("Expected QMatrix is different from Resulting QMatrix %v", sub.At(i, j))
			}
		}
	}
}

func TestTN93Eigens(t *testing.T) {
	var v, l, r []float64
	var vMat, lMat, rMat mat.Matrix
	var err error
	var expQMatrix *mat.Dense
	var resQMatrix *mat.Dense = mat.NewDense(4, 4, nil)
	var sub *mat.Dense = mat.NewDense(4, 4, nil)

	expQMatrix = mat.NewDense(4, 4, []float64{
		-1.0000000, 0.3333333, 0.3333333, 0.3333333,
		0.3333333, -1.0000000, 0.3333333, 0.3333333,
		0.3333333, 0.3333333, -1.0000000, 0.3333333,
		0.3333333, 0.3333333, 0.3333333, -1.0000000,
	})
	t.Logf("expQ = %v", mat.Formatted(expQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	m := NewTN93Model()
	m.InitModel(1., 1., 1./4., 1./4., 1./4., 1./4.)

	if v, l, r, err = m.Eigens(); err != nil {
		t.Errorf("Error while computing F81 eigen vectors: %v", err)
	}

	// Transpose because it was in col-major formrat
	lMat = mat.NewDense(4, 4, l).T()
	rMat = mat.NewDense(4, 4, r).T()

	vMat = mat.NewDiagDense(4, v)

	resQMatrix.Mul(lMat, vMat)
	resQMatrix.Mul(resQMatrix, rMat)
	sub.Sub(resQMatrix, expQMatrix)
	t.Logf("resQ = %v", mat.Formatted(resQMatrix, mat.Prefix("                 "), mat.Squeeze()))

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if sub.At(i, j) > .0000001 {
				t.Errorf("Expected QMatrix is different from Resulting QMatrix %v", sub.At(i, j))
			}
		}
	}
}

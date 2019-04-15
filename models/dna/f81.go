package dna

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type F81Model struct {
	// Parameters (for eigen values/vectors computation)
	// See https://en.wikipedia.org/wiki/Models_of_DNA_evolution#F81_model_(Felsenstein_1981)
	qmatrix *mat.Dense
}

func NewF81Model() *F81Model {
	return &F81Model{
		nil,
	}
}

func (m *F81Model) InitModel(piA, piC, piG, piT float64) {
	m.qmatrix = mat.NewDense(4, 4, []float64{
		-(piC + piG + piT), piC, piG, piT,
		piA, -(piA + piG + piT), piG, piT,
		piA, piC, -(piA + piC + piT), piT,
		piA, piC, piG, -(piA + piC + piG),
	})
	// Normalization of Q
	norm := 1. / (2 * (piA*piC + piA*piG + piA*piT + piC*piG + piC*piT + piG*piT))
	m.qmatrix.Apply(func(i, j int, v float64) float64 { return v * norm }, m.qmatrix)
}

func (m *F81Model) Eigens() (val []float64, leftvectors, rightvectors []float64, err error) {
	// Compute eigen values, left and right eigenvectors of Q
	eigen := &mat.Eigen{}
	if ok := eigen.Factorize(m.qmatrix, mat.EigenRight); !ok {
		err = fmt.Errorf("Problem during matrix decomposition")
		return
	}

	val = make([]float64, 4)
	for i, b := range eigen.Values(nil) {
		val[i] = real(b)
	}
	u := eigen.VectorsTo(nil)
	reigenvect := mat.NewDense(4, 4, nil)
	leigenvect := mat.NewDense(4, 4, nil)
	reigenvect.Apply(func(i, j int, val float64) float64 { return real(u.At(i, j)) }, reigenvect)
	leigenvect.Inverse(reigenvect)

	leftvectors = leigenvect.RawMatrix().Data
	rightvectors = reigenvect.RawMatrix().Data

	return
}

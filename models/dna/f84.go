package dna

type F84Model struct {
	// Parameters (for eigen values/vectors computation)
	// https://en.wikipedia.org/wiki/Models_of_DNA_evolution#HKY85_model_(Hasegawa,_Kishino_and_Yano_1985)
	piA, piC, piG, piT float64
	kappa              float64
}

func NewF84Model() *F84Model {
	return &F84Model{
		1. / 4., 1. / 4., 1. / 4., 1. / 4.,
		1.0,
	}
}

func (m *F84Model) InitModel(kappa, piA, piC, piG, piT float64) {
	//m.qmatrix = mat.NewDense(4, 4, []float64{
	//	-(piC + (1+kappa/piR)*piG + piT), piC, (1 + kappa/piR) * piG, piT,
	//	piA, -(piA + piG + (1+kappa/piY)*piT), piG, (1 + kappa/piY) * piT,
	//	(1 + kappa/piR) * piA, piC, -((1+kappa/piR)*piA + piC + piT), piT,
	//	piA, (1 + kappa/piY) * piC, piG, -(piA + (1+kappa/piY)*piC + piG),
	//})
	// Normalization of Q
	m.kappa = kappa
	m.piA = piA
	m.piC = piC
	m.piG = piG
	m.piT = piT
}

// See http://biopp.univ-montp2.fr/Documents/ClassDocumentation/bpp-phyl/html/F84_8cpp_source.html
func (m *F84Model) Eigens() (val []float64, leftvectors, rightvectors []float64, err error) {
	piY := m.piT + m.piC
	piR := m.piA + m.piG
	norm := 1. / (1 - m.piA*m.piA - m.piC*m.piC - m.piG*m.piG - m.piT*m.piT + 2.*m.kappa*(m.piC*m.piT/piY+m.piA*m.piG/piR))

	val = []float64{
		0,
		-norm * (1 + m.kappa),
		-norm * (1 + m.kappa),
		-norm,
	}

	leftvectors = []float64{
		m.piA, m.piC, m.piG, m.piT,
		0., m.piT / piY, 0., -m.piT / piY,
		m.piG / piR, 0., -m.piG / piR, 0.,
		m.piA * piY / piR, -m.piC, m.piG * piY / piR, -m.piT,
	}

	rightvectors = []float64{
		1., 0., 1., 1.,
		1., 1., 0., -piR / piY,
		1., 0., -m.piA / m.piG, 1.,
		1., -m.piC / m.piT, 0., -piR / piY,
	}

	return
}

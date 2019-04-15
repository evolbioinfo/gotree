package dna

type JCModel struct {
}

func NewJCModel() *JCModel {
	return &JCModel{}
}

func (m *JCModel) InitModel() (err error) {
	return
}

// Left vectors and right vectors are given in column-major format
func (m *JCModel) Eigens() (val []float64, leftvectors, rightvectors []float64, err error) {
	val = []float64{
		0,
		-4. / 3.,
		-4. / 3.,
		-4. / 3.,
	}

	leftvectors = []float64{
		1. / 4., 1. / 4., 1. / 4., 1. / 4.,
		-1. / 4., -1. / 4., 3. / 4., -1. / 4.,
		-1. / 4., 3. / 4., -1. / 4., -1. / 4.,
		3. / 4., -1. / 4., -1. / 4., -1. / 4.,
	}

	rightvectors = []float64{
		1., 0., 0., 1.,
		1., 0., 1., 0.,
		1., 1., 0., 0.,
		1., -1., -1., -1.,
	}
	return
}

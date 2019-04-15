package dna

type DNAModel interface {
	Eigens() (val []float64, leftvectors, rightvectors []float64, err error)
}

// package acr provides data structures & functions
// for ancestral sequence reconstruction methods
package asr

import (
	"errors"
	"fmt"
)

type AncestralState struct {
	counts []float64 // counts or probabilities of each character indices correspond to goalign/align/AlphabetCharacters()
}

type AncestralSequence struct {
	seq []AncestralState // Sequence of ancestral states
}

func NewAncestralSequence(length int, alphabetlength int) (*AncestralSequence, error) {
	if length < 0 || alphabetlength < 0 {
		return nil, errors.New(fmt.Sprintf("Cannot initialize an Ancestral Sequence of length %d with an alphabet of length %d", length, alphabetlength))
	}
	seq := &AncestralSequence{make([]AncestralState, length)}
	for i := 0; i < length; i++ {
		seq.seq[i] = AncestralState{make([]float64, alphabetlength)}
	}
	return seq, nil
}

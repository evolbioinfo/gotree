// package acr provides data structures & functions
// for ancestral character reconstruction methods
package acr

// counts or probabilities of each character
// indices correspond to Alphanumeric order of the alphabet
// (all possible states)
type AncestralState []float64

func AncestralStateIndices(alphabet []string) map[string]int {
	indices := make(map[string]int)
	for i, v := range alphabet {
		indices[v] = i
	}
	return indices
}

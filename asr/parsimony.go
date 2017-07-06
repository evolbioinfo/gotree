package asr

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/fredericlemoine/goalign/align"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/tree"
)

// Will annotate the tree nodes with ancestral sequences
// Computed using parsimony
// Sequences will be located in the comment field of each node
// at the first index
func ParsimonyAsr(t *tree.Tree, a align.Alignment) error {
	var err error
	var nodes []*tree.Node = t.Nodes()
	var seqs []*AncestralSequence = make([]*AncestralSequence, len(nodes))
	// Initialize indices of characters
	var charToIndex map[rune]int = make(map[rune]int)
	for i, c := range a.AlphabetCharacters() {
		charToIndex[c] = i
	}

	// We initialize all ancestral sequences
	// And sequences at tips
	for i, n := range nodes {
		n.SetId(i)
		if seqs[i], err = NewAncestralSequence(a.Length(), len(a.AlphabetCharacters())); err != nil {
			return err
		}
	}

	err = parsimonyPostOrder(t.Root(), nil, a, seqs, charToIndex)
	if err != nil {
		return err
	}
	parsimonyPreOrder(t.Root(), nil, a, seqs, charToIndex)
	assignSequencesToTree(t, seqs, a.AlphabetCharacters())
	return nil
}

// First step of the parsimony computatation: From tips to root
func parsimonyPostOrder(cur, prev *tree.Node, a align.Alignment, seqs []*AncestralSequence, charToIndex map[rune]int) error {
	// If it is a tip, we initialize the ancestral sequences using the current
	// Sequence in the alignment. If no such sequence exists in the alignment,
	// then returns an error
	if cur.Tip() {
		seq, ok := a.GetSequenceChar(cur.Name())
		if !ok {
			return errors.New(fmt.Sprintf("Sequence %s does not exist in the alignment", cur.Name()))
		}
		for j, c := range seq {
			charindex, ok := charToIndex[c]
			if ok {
				seqs[cur.Id()].seq[j].counts[charindex] = 1
			} else {
				io.LogWarning(errors.New(fmt.Sprintf("Character %c does not exist in the alphabet, ignoring the state", c)))
			}
		}
	} else {
		for _, child := range cur.Neigh() {
			if child != prev {
				if err := parsimonyPostOrder(child, cur, a, seqs, charToIndex); err != nil {
					return err
				}
			}
		}
		// As we are manipulating trees with multifurcations
		// For each character we count the number of children having it
		// and then we take character(s) with the maximum number of children
		// And that for each site of the alignment
		for j, ances := range seqs[cur.Id()].seq {
			for _, child := range cur.Neigh() {
				if child != prev {
					counts := seqs[child.Id()].seq[j].counts
					for k, c := range counts {
						ances.counts[k] += c
					}
				}
			}
			// Now we set to 0 all character states that are not the max, and to 1 the states that are the max
			max := 0.0
			for _, c := range ances.counts {
				if c > max {
					max = c
				}
			}
			for k, c := range ances.counts {
				if c == max {
					ances.counts[k] = 1
				} else {
					ances.counts[k] = 0
				}
			}
		}
	}
	return nil
}

// Second step of the parsimony computatation: From root to tips
func parsimonyPreOrder(cur, prev *tree.Node, a align.Alignment, seqs []*AncestralSequence, charToIndex map[rune]int) {
	// If it is not a tip and not the root
	if !cur.Tip() {
		if prev != nil {
			// As we are manipulating trees with multifurcations
			// For each character we count the number of children having it
			// and then we take character(s) with the maximum number of children
			// And that for each site of the alignment
			for j, ances := range seqs[cur.Id()].seq {
				state := AncestralState{make([]float64, len(charToIndex))}
				// With Parent
				for _, child := range cur.Neigh() {
					counts := seqs[child.Id()].seq[j].counts
					for k, c := range counts {
						state.counts[k] += c
					}
				}
				// Now we set to 0 all character states that are not the max, and to 1 the states that are the max
				maxall := 0.0
				nbmaxall := 0
				for _, c := range state.counts {
					if c > maxall {
						maxall = c
						nbmaxall = 1
					} else if c == maxall {
						nbmaxall++
					}
				}
				for k, c := range state.counts {
					if c == maxall {
						ances.counts[k] = 1
					} else {
						ances.counts[k] = 0
					}
				}
			}
		}
		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyPreOrder(child, cur, a, seqs, charToIndex)
			}
		}
	}
}
func assignSequencesToTree(t *tree.Tree, seqs []*AncestralSequence, alphabet []rune) {
	var buffer bytes.Buffer
	var subbuffer bytes.Buffer

	for _, n := range t.Nodes() {
		buffer.Reset()
		ancseq := seqs[n.Id()]
		for _, state := range ancseq.seq {
			subbuffer.Reset()
			nb := 0
			for i, c := range state.counts {
				if c > 0 {
					subbuffer.WriteRune(alphabet[i])
					nb++
				}
			}
			// If no state has a count> 0 : All are possible
			// -
			if nb == 0 {
				subbuffer.WriteRune('*')
			}
			if nb > 1 {
				buffer.WriteRune('{')
			}
			buffer.Write(subbuffer.Bytes())
			if nb > 1 {
				buffer.WriteRune('}')
			}
		}
		n.ClearComments()
		n.AddComment(buffer.String())
	}
}

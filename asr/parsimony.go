package asr

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/evolbioinfo/goalign/align"
	"github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/tree"
)

const (
	ALGO_DELTRAN = iota
	ALGO_ACCTRAN
	ALGO_DOWNPASS
	ALGO_NONE
)

// Will annotate the tree nodes with ancestral sequences
// Computed using parsimony
// Sequences will be located in the comment field of each node
// at the first index
func ParsimonyAsr(t *tree.Tree, a align.Alignment, algo int, randomResolve bool) (nsteps []int, err error) {
	var nodes []*tree.Node = t.Nodes()
	var seqs []*AncestralSequence = make([]*AncestralSequence, len(nodes))
	var upseqs []*AncestralSequence = make([]*AncestralSequence, len(nodes)) // Upside seqs of each  node
	var alphabet []uint8 = a.AlphabetCharacters()

	alphabet = append(alphabet, '-')
	alphabet = append(alphabet, '*')

	// Initialize indices of characters
	var charToIndex map[uint8]int = make(map[uint8]int)
	for i, c := range alphabet {
		charToIndex[c] = i
	}
	nsteps = make([]int, a.Length()+1)

	// We initialize all ancestral sequences
	// And sequences at tips
	for i, n := range nodes {
		n.SetId(i)
		if seqs[i], err = NewAncestralSequence(a.Length(), len(charToIndex)); err != nil {
			return nil, err
		}
		if upseqs[i], err = NewAncestralSequence(a.Length(), len(charToIndex)); err != nil {
			return nil, err
		}
	}

	err = parsimonyUPPASS(t.Root(), nil, a, seqs, nsteps, charToIndex)
	if err != nil {
		return
	}

	switch algo {
	case ALGO_DOWNPASS:
		parsimonyDOWNPASS(t.Root(), nil, a, seqs, upseqs, charToIndex, randomResolve)
	case ALGO_DELTRAN:
		parsimonyDOWNPASS(t.Root(), nil, a, seqs, upseqs, charToIndex, false)
		parsimonyDELTRAN(t.Root(), nil, a, seqs, charToIndex, randomResolve)
	case ALGO_ACCTRAN:
		parsimonyACCTRAN(t.Root(), nil, a, seqs, charToIndex, randomResolve)
	default:
		err = fmt.Errorf("parsimony algorithm %d unkown", algo)
		return
	}

	assignSequencesToTree(t, seqs, alphabet)
	return
}

// First step of the parsimony computatation: From tips to root
func parsimonyUPPASS(cur, prev *tree.Node, a align.Alignment, seqs []*AncestralSequence, nsteps []int, charToIndex map[uint8]int) (err error) {
	// If it is a tip, we initialize the ancestral sequences using the current
	// Sequence in the alignment. If no such sequence exists in the alignment,
	// then returns an error
	if cur.Tip() {
		seq, ok := a.GetSequenceChar(cur.Name())
		if !ok {
			err = fmt.Errorf("sequence %s does not exist in the alignment", cur.Name())
			return
		}
		for j, c := range seq {
			possibilities := make([]uint8, 0)
			if a.Alphabet() == align.NUCLEOTIDS {
				possibilities = align.IupacCode[c]
			} else {
				if c == align.ALL_AMINO {
					for k := range charToIndex {
						possibilities = append(possibilities, k)
					}
					possibilities = possibilities[:len(possibilities)-2]
				} else {
					possibilities = append(possibilities, c)
				}
			}
			for _, c2 := range possibilities {
				charindex, ok := charToIndex[c2]
				if ok {
					seqs[cur.Id()].seq[j].counts[charindex] = 1
				} else {
					io.LogWarning(fmt.Errorf("character %c does not exist in the alphabet, ignoring the state", c2))
				}
			}
		}
	} else {
		for _, child := range cur.Neigh() {
			if child != prev {
				if err = parsimonyUPPASS(child, cur, a, seqs, nsteps, charToIndex); err != nil {
					return
				}
			}
		}
		// As we are manipulating trees with multifurcations
		// For each character we count the number of children having it
		// and then we take character(s) with the maximum number of children
		// And that for each site of the alignment
		for j, ances := range seqs[cur.Id()].seq {
			nchild := 0
			for _, child := range cur.Neigh() {
				if child != prev {
					counts := seqs[child.Id()].seq[j].counts
					for k, c := range counts {
						ances.counts[k] += c
					}
					nchild++
				}
			}
			maxState := 0
			max := 0.0
			for k, c := range seqs[cur.Id()].seq[j].counts {
				if c > max {
					maxState = k
					max = c
				}
			}
			computeParsimony(ances, ances, nchild)
			for _, child := range cur.Neigh() {
				if child != prev {
					if seqs[child.Id()].seq[j].counts[maxState] == 0 {
						nsteps[j]++
					}
				}
			}
		}
	}
	return
}

// Second step of the parsimony computatation: From root to tips
func parsimonyDOWNPASS(cur, prev *tree.Node, a align.Alignment,
	seqs []*AncestralSequence, upseqs []*AncestralSequence,
	charToIndex map[uint8]int, randomResolve bool) {
	// If it is not a tip and not the root
	if !cur.Tip() {
		// We compute the up sequence states for each children of
		// the current node (may be the root)
		// i.e. the parsimony from the upside of the tree
		for _, child := range cur.Neigh() {
			if child != prev {
				for j, _ := range seqs[cur.Id()].seq {
					state := AncestralState{make([]float64, len(charToIndex))}
					nchild := 0
					// already computed up state of the current node
					if prev != nil { // Not the root
						nchild++
						for k, c := range upseqs[cur.Id()].seq[j].counts {
							state.counts[k] += c
						}
					}
					// already computed down states of children of current node
					// except current child _child_
					for _, child2 := range cur.Neigh() {
						if child2 != prev && child2 != child {
							for k, c := range seqs[child2.Id()].seq[j].counts {
								state.counts[k] += c
							}
							nchild++
						}
					}
					// Compute the up state now
					computeParsimony(state, upseqs[child.Id()].seq[j], nchild)
				}
			}
		}

		if prev != nil {
			// As we are manipulating trees with multifurcations
			// For each character we count the number of children having it
			// and then we take character(s) with the maximum number of children
			// And that for each site of the alignment
			for j, _ := range seqs[cur.Id()].seq {
				state := AncestralState{make([]float64, len(charToIndex))}
				// With Parent using its upseq
				nchild := 1
				for k, c := range upseqs[cur.Id()].seq[j].counts {
					state.counts[k] += c
				}
				for _, child := range cur.Neigh() {
					if child != prev {
						for k, c := range seqs[child.Id()].seq[j].counts {
							state.counts[k] += c
						}
						nchild++
					}
				}
				computeParsimony(state, seqs[cur.Id()].seq[j], nchild)
			}
		}

		// We randomly resolve ambiguities
		// Even for the root (outside if statement)
		if randomResolve {
			randomlyResolveNodeStates(cur, seqs)
		}

		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyDOWNPASS(child, cur, a, seqs, upseqs, charToIndex, randomResolve)
			}
		}
	}
}

func computeParsimony(neighborStates AncestralState, currentStates AncestralState, nchild int) {
	// If intersection of states of children and parent is emtpy:
	// then State of cur node ==  Union of intersection of nodes 2 by 2
	// If state is still empty, then state of cur node is union of all states
	max := 0.0
	for _, c := range neighborStates.counts {
		if c > max {
			max = c
		}
	}
	for k, c := range neighborStates.counts {
		if c == max {
			// We have a characters shared by all neighbors and parent: Intersection ok
			currentStates.counts[k] = 1
		} else {
			// Else we do not take it
			currentStates.counts[k] = 0
		}
	}
}

// Third step of the parsimony computation for resolving ambiguities
func parsimonyDELTRAN(cur, prev *tree.Node, a align.Alignment, seqs []*AncestralSequence, charToIndex map[uint8]int, randomResolve bool) {
	// If it is not a tip
	if !cur.Tip() {
		// If it is not the root
		if prev != nil {
			for j, ances := range seqs[cur.Id()].seq {
				state := AncestralState{make([]float64, len(charToIndex))}
				// Compute the intersection with Parent
				nullIntersection := true
				for k, c := range ances.counts {
					state.counts[k] += c
				}
				for k, c := range seqs[prev.Id()].seq[j].counts {
					state.counts[k] += c
					if state.counts[k] > 1 {
						nullIntersection = false
					}
				}

				// If non null intersection, then current node's state is the intersection
				if !nullIntersection {
					for k, c := range state.counts {
						if c > 1 {
							ances.counts[k] = 1
						} else {
							ances.counts[k] = 0
						}
					}
				}
			}
		}

		// We resolve ambiguities if randomResolve
		// Even for the root (outside if statement)
		if randomResolve {
			randomlyResolveNodeStates(cur, seqs)
		}

		// We go down in the tree
		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyDELTRAN(child, cur, a, seqs, charToIndex, randomResolve)
			}
		}
	}
}

// Second step of the parsimony computation (instead of DOWNPASS) for resolving ambiguities
func parsimonyACCTRAN(cur, prev *tree.Node, a align.Alignment, seqs []*AncestralSequence, charToIndex map[uint8]int, randomResolve bool) {
	// If it is not a tip
	if !cur.Tip() {
		// We resolve ambiguities if randomResolve
		if randomResolve {
			randomlyResolveNodeStates(cur, seqs)
		}

		// We Analyze each direct child
		for _, child := range cur.Neigh() {
			if child != prev {
				for j, ances := range seqs[cur.Id()].seq {
					state := AncestralState{make([]float64, len(charToIndex))}
					// Compute the intersection with Parent
					nullIntersection := true
					for k, c := range seqs[child.Id()].seq[j].counts {
						state.counts[k] += c
					}
					for k, c := range ances.counts {
						state.counts[k] += c
						if state.counts[k] > 1 {
							nullIntersection = false
						}
					}
					// If non null intersection, then child node's state is the intersection
					if !nullIntersection {
						for k, c := range state.counts {
							if c > 1 {
								seqs[child.Id()].seq[j].counts[k] = 1
							} else {
								seqs[child.Id()].seq[j].counts[k] = 0
							}
						}
					}
				}
			}
		}
		// We go down in the tree
		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyACCTRAN(child, cur, a, seqs, charToIndex, randomResolve)
			}
		}
	}
}

// Randomly resolve all sequences states of the node that are ambiguous
func randomlyResolveNodeStates(node *tree.Node, seqs []*AncestralSequence) {
	for _, ances := range seqs[node.Id()].seq {
		numstates := 0
		for _, c := range ances.counts {
			if c >= 1 {
				numstates++
			}
		}
		if numstates > 1 {
			curstate := 0
			randstate := rand.Intn(numstates)
			for k, c := range ances.counts {
				if c >= 1 {
					if curstate == randstate {
						ances.counts[k] = 1
					} else {
						ances.counts[k] = 0
					}
					curstate++
				} else {
					ances.counts[k] = 0
				}
			}
		}
	}
}

func assignSequencesToTree(t *tree.Tree, seqs []*AncestralSequence, alphabet []uint8) {
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
					subbuffer.WriteByte(alphabet[i])
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
		n.AddComment(buffer.String())
	}
}

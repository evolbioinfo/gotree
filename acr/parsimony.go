package acr

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/fredericlemoine/gotree/tree"
)

// Will annotate the tree nodes with ancestral characters
// Computed using parsimony
// Characters will be located in the comment field of each node
// at the first index
// tipCharacters: mapping between tipnames and character state
func ParsimonyAcr(t *tree.Tree, tipCharacters map[string]string) error {
	var err error
	var nodes []*tree.Node = t.Nodes()
	var states []AncestralState = make([]AncestralState, len(nodes))
	// Initialize indices of characters
	alphabet := make([]string, 0, 10)
	sort.Strings(alphabet)
	seenState := make(map[string]bool)
	for _, state := range tipCharacters {
		if _, ok := seenState[state]; !ok {
			alphabet = append(alphabet, state)
		}
		seenState[state] = true
	}
	stateIndices := AncestralStateIndices(alphabet)

	// We initialize all ancestral states
	for i, n := range nodes {
		n.SetId(i)
		states[i] = make(AncestralState, len(alphabet))
	}

	err = parsimonyPostOrder(t.Root(), nil, tipCharacters, states, stateIndices)
	if err != nil {
		return err
	}
	parsimonyPreOrder(t.Root(), nil, states, stateIndices)
	assignStatesToTree(t, states, alphabet)
	return nil
}

// First step of the parsimony computatation: From tips to root
func parsimonyPostOrder(cur, prev *tree.Node, tipCharacters map[string]string, states []AncestralState, stateIndices map[string]int) error {
	// If it is a tip, we initialize the ancestral state using the current
	// state in the alignment. If no such tip name exists in the mapping file,
	// then returns an error
	if cur.Tip() {
		state, ok := tipCharacters[cur.Name()]
		if !ok {
			return errors.New(fmt.Sprintf("Tip %s does not exist in the tip/state mapping file", cur.Name()))
		}
		stateindex, ok := stateIndices[state]
		if ok {
			states[cur.Id()][stateindex] = 1
		} else {
			return errors.New(fmt.Sprintf("State %c does not exist in the alphabet, ignoring the state", state))
		}
	} else {
		for _, child := range cur.Neigh() {
			if child != prev {
				if err := parsimonyPostOrder(child, cur, tipCharacters, states, stateIndices); err != nil {
					return err
				}
			}
		}
		// As we are manipulating trees with multifurcations
		// For each character we count the number of children having it
		// and then we take character(s) with the maximum number of children
		for _, child := range cur.Neigh() {
			if child != prev {
				for k, c := range states[child.Id()] {
					states[cur.Id()][k] += c
				}
			}
		}
		// Now we set to 0 all character states that are not the max, and to 1 the states that are the max
		max := 0.0
		for _, c := range states[cur.Id()] {
			if c > max {
				max = c
			}
		}
		for k, c := range states[cur.Id()] {
			if c == max {
				states[cur.Id()][k] = 1
			} else {
				states[cur.Id()][k] = 0
			}
		}
	}
	return nil
}

// Second step of the parsimony computatation: From root to tips
func parsimonyPreOrder(cur, prev *tree.Node, states []AncestralState, stateIndices map[string]int) {
	// If it is not a tip and not the root
	if !cur.Tip() {
		if prev != nil {
			// As we are manipulating trees with multifurcations
			// For each character we count the number of children having it
			// and then we take character(s) with the maximum number of children
			state := make(AncestralState, len(stateIndices))
			// With Parent
			for _, child := range cur.Neigh() {
				for k, c := range states[child.Id()] {
					state[k] += c
				}
			}
			// Now we set to 0 all character states that are not the max, and to 1 the states that are the max
			maxall := 0.0
			nbmaxall := 0
			for _, c := range state {
				if c > maxall {
					maxall = c
					nbmaxall = 1
				} else if c == maxall {
					nbmaxall++
				}
			}
			for k, c := range state {
				if c == maxall {
					states[cur.Id()][k] = 1
				} else {
					states[cur.Id()][k] = 0
				}
			}
		}

		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyPreOrder(child, cur, states, stateIndices)
			}
		}
	}
}

func assignStatesToTree(t *tree.Tree, states []AncestralState, alphabet []string) {
	var buffer bytes.Buffer

	for _, n := range t.Nodes() {
		buffer.Reset()
		nb := 0
		for i, c := range states[n.Id()] {
			if c > 0 {
				if nb > 0 {
					buffer.WriteRune('|')
				}
				buffer.WriteString(alphabet[i])
				nb++
			}
		}
		// If no state has a count> 0 : All are possible
		// *
		if nb == 0 {
			buffer.WriteRune('*')
		}
		n.ClearComments()
		n.AddComment(buffer.String())
	}
}

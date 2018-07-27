package acr

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/fredericlemoine/gotree/tree"
)

const (
	ALGO_DELTRAN = iota
	ALGO_ACCTRAN
	ALGO_DOWNPASS
	ALGO_NONE
)

// Will annotate the tree nodes with ancestral characters
// Computed using parsimony
// Characters will be located in the comment field of each node
// at the first index
// tipCharacters: mapping between tipnames and character state
// Algo: One of ALGO_DELTRAN, ALGO_ACCTRAN, and ALGO_DOWNPASS : returns an error otherwise
// If ALGO_DOWNPASS, just executes UPPASS then DOWNPASS,
// If ALGO_DELTRAN, then executes UPPASS, DOWNPASS, then DELTRAN,
// If ALGO_ACCTRAN, then executes UPPASS and ACCTRAN
// Returns a map with the states of all nodes. If a node has a name, key is its name, if a node has no name,
// the key will be its id in the deep first traversal of the tree.
// if randomResolve is true, then in the second pass, each ambiguities will be resolved randomly
func ParsimonyAcr(t *tree.Tree, tipCharacters map[string]string, algo int, randomResolve bool) (map[string]string, error) {
	var err error
	var nodes []*tree.Node = t.Nodes()
	var states []AncestralState = make([]AncestralState, len(nodes))   // Downside states of each node
	var upstates []AncestralState = make([]AncestralState, len(nodes)) // Upside states of each  node
	// Initialize indices of characters
	alphabet := make([]string, 0, 10)
	seenState := make(map[string]bool)
	for _, state := range tipCharacters {
		if _, ok := seenState[state]; !ok {
			alphabet = append(alphabet, state)
		}
		seenState[state] = true
	}
	sort.Strings(alphabet)
	stateIndices := AncestralStateIndices(alphabet)

	// We initialize all ancestral states
	for i, n := range nodes {
		n.SetId(i)
		states[i] = make(AncestralState, len(alphabet))
		upstates[i] = make(AncestralState, len(alphabet))
	}

	err = parsimonyUPPASS(t.Root(), nil, tipCharacters, states, stateIndices)
	if err != nil {
		return nil, err
	}

	switch algo {
	case ALGO_DOWNPASS:
		parsimonyDOWNPASS(t.Root(), nil, states, upstates, stateIndices, randomResolve)
	case ALGO_DELTRAN:
		parsimonyDOWNPASS(t.Root(), nil, states, upstates, stateIndices, false)
		parsimonyDELTRAN(t.Root(), nil, states, stateIndices, randomResolve)
	case ALGO_ACCTRAN:
		parsimonyACCTRAN(t.Root(), nil, states, stateIndices, randomResolve)
	case ALGO_NONE:
		// No pass after uppass
	default:
		return nil, fmt.Errorf("Parsimony algorithm %d unkown", algo)
	}

	nametostates := buildInternalNamesToStatesMap(t, states, alphabet)
	assignStatesToTree(t, states, alphabet)
	return nametostates, nil
}

// First step of the parsimony computatation: From tips to root
func parsimonyUPPASS(cur, prev *tree.Node, tipCharacters map[string]string, states []AncestralState, stateIndices map[string]int) error {
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
			return errors.New(fmt.Sprintf("State %s does not exist in the alphabet, ignoring the state", state))
		}
	} else {
		for _, child := range cur.Neigh() {
			if child != prev {
				if err := parsimonyUPPASS(child, cur, tipCharacters, states, stateIndices); err != nil {
					return err
				}
			}
		}

		// If intersection of states of children is emtpy:
		// then State of cur node ==  Union of State of children if
		// Else
		// State of cur node ==  Intersection of States of children if
		// works with trees having multifurcations
		nchild := 0
		for _, child := range cur.Neigh() {
			if child != prev {
				for k, c := range states[child.Id()] {
					states[cur.Id()][k] += c
				}
				nchild++
			}
		}
		computeParsimony(states[cur.Id()], states[cur.Id()], nchild)
	}
	return nil
}

// Second step of the parsimony computation: From root to tips
func parsimonyDOWNPASS(cur, prev *tree.Node,
	states []AncestralState, upstates []AncestralState,
	stateIndices map[string]int, randomResolve bool) {
	// If it is not a tip and not the root
	if !cur.Tip() {
		// We compute the up state for each children of
		// the current node (may be the root)
		// i.e. the parsimony from the upside of the tree
		for _, child := range cur.Neigh() {
			if child != prev {
				state := make(AncestralState, len(stateIndices))
				nchild := 0
				// already computed up state of the current node
				if prev != nil { // Not the root
					nchild++
					for k, c := range upstates[cur.Id()] {
						state[k] += c
					}
				}
				// already computed down states of children of current node
				// except current child _child_
				for _, child2 := range cur.Neigh() {
					if child2 != prev && child2 != child {
						for k, c := range states[child2.Id()] {
							state[k] += c
						}
						nchild++
					}
				}
				// Compute the up state now
				computeParsimony(state, upstates[child.Id()], nchild)
			}
		}

		// Not the root
		if prev != nil {
			// As we are manipulating trees with multifurcations
			// For each character we count the number of children having it
			// and then we take character(s) with the maximum number of children
			state := make(AncestralState, len(stateIndices))
			// With Parent (upstate of cur node)
			nchild := 1
			for k, c := range upstates[cur.Id()] {
				state[k] += c
			}
			for _, child := range cur.Neigh() {
				if child != prev {
					for k, c := range states[child.Id()] {
						state[k] += c
					}
					nchild++
				}
			}
			computeParsimony(state, states[cur.Id()], nchild)
		}
		// We randomly resolve ambiguities
		// Even for the root (outside if statement)
		if randomResolve {
			randomlyResolveNodeStates(cur, states)
		}

		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyDOWNPASS(child, cur, states, upstates, stateIndices, randomResolve)
			}
		}
	}
}

// Will set the most parsimonious states in the "currentStates" slice
// using the neighbor states "neighborStates", and the number of neighbors
func computeParsimony(neighborStates AncestralState, currentStates AncestralState, nchild int) {
	// If intersection of states of children and parent is emtpy:
	// then State of cur node ==  Union of intersection of nodes 2 by 2
	// If state is still empty, then state of cur node is union of all states
	max := 0.0
	for _, c := range neighborStates {
		if c > max {
			max = c
		}
	}
	for k, c := range neighborStates {
		if int(max) == nchild && c == max {
			// We have a characters shared by all neighbors and parent: Intersection ok
			currentStates[k] = 1
		} else if int(max) == 1 && c > 0 {
			// Else we have no intersection between any children: take union
			currentStates[k] = 1
		} else if int(max) < nchild && c > 1 {
			// Else we have a character shared by at least 2 children: OK
			currentStates[k] = 1
		} else {
			// Else we do not take it
			currentStates[k] = 0
		}
	}
}

// Third step of the parsimony computation for resolving ambiguities
func parsimonyDELTRAN(cur, prev *tree.Node, states []AncestralState, stateIndices map[string]int, randomResolve bool) {
	// If it is not a tip
	if !cur.Tip() {
		// If it is not the root
		if prev != nil {
			state := make(AncestralState, len(stateIndices))
			// Compute the intersection with Parent
			nullIntersection := true
			for k, c := range states[cur.Id()] {
				state[k] += c
			}
			for k, c := range states[prev.Id()] {
				state[k] += c
				if state[k] > 1 {
					nullIntersection = false
				}
			}

			// If non null intersection, then current node's state is the intersection
			if !nullIntersection {
				for k, c := range state {
					if c > 1 {
						states[cur.Id()][k] = 1
					} else {
						states[cur.Id()][k] = 0
					}
				}
			}
		}
		// We resolve ambiguities if randomResolve
		// Even for the root (outside if statement)
		if randomResolve {
			randomlyResolveNodeStates(cur, states)
		}

		// We go down in the tree
		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyDELTRAN(child, cur, states, stateIndices, randomResolve)
			}
		}
	}
}

// Second step of the parsimony computation (instead of DOWNPASS) for resolving ambiguities
func parsimonyACCTRAN(cur, prev *tree.Node, states []AncestralState, stateIndices map[string]int, randomResolve bool) {
	// If it is not a tip
	if !cur.Tip() {
		// We resolve the root ambiguities if randomResolve
		if randomResolve {
			randomlyResolveNodeStates(cur, states)
		}

		// We Analyze each direct child
		for _, child := range cur.Neigh() {
			if child != prev {
				state := make(AncestralState, len(stateIndices))
				// Compute the intersection with Parent
				nullIntersection := true
				for k, c := range states[child.Id()] {
					state[k] += c
				}
				for k, c := range states[cur.Id()] {
					state[k] += c
					if state[k] > 1 {
						nullIntersection = false
					}
				}

				// If non null intersection, then child node's state is the intersection
				if !nullIntersection {
					for k, c := range state {
						if c > 1 {
							states[child.Id()][k] = 1
						} else {
							states[child.Id()][k] = 0
						}
					}
				}
			}
		}
		// We go down in the tree
		for _, child := range cur.Neigh() {
			if child != prev {
				parsimonyACCTRAN(child, cur, states, stateIndices, randomResolve)
			}
		}
	}
}

func randomlyResolveNodeStates(node *tree.Node, states []AncestralState) {
	numstates := 0
	for _, c := range states[node.Id()] {
		if c >= 1 {
			numstates++
		}
	}
	if numstates > 1 {
		curstate := 0
		randstate := rand.Intn(numstates)
		for k, c := range states[node.Id()] {
			if c >= 1 {
				if curstate == randstate {
					states[node.Id()][k] = 1
				} else {
					states[node.Id()][k] = 0
				}
				curstate++
			} else {
				states[node.Id()][k] = 0
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

// Returns a map with keys: Internal nodes identifier (id or name if any), and value: list of possible states, comma separated
func buildInternalNamesToStatesMap(t *tree.Tree, states []AncestralState, alphabet []string) map[string]string {
	outmap := make(map[string]string)
	st := make([]string, 0, 10)

	for _, n := range t.Nodes() {
		if !n.Tip() {
			nb := 0
			st = st[:0]
			for i, c := range states[n.Id()] {
				if c > 0 {
					st = append(st, alphabet[i])
					nb++
				}
			}
			// If no state has a count> 0 : All are possible
			// *
			if nb == 0 {
				st = append(st, "*")
			}

			id := fmt.Sprintf("%d", n.Id())
			if n.Name() != "" {
				id = n.Name()
			}
			sort.Strings(st)
			outmap[id] = strings.Join(st, ",")
		}
	}
	return outmap
}

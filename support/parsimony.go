package support

import (
	"github.com/fredericlemoine/gostats"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"math"
	"sort"
	"sync"
)

const (
	STATENULL = -1 // Null state
	STATE0    = 0  // State 0
	STATE1    = 1  // State 1
	STATE01   = 2  // State 0/1
)

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

type parsimonySupporter struct {
	distribution_rand_step_val [][]float64
	expected_rand_steps        []float64
}

// This function precomputes the esperence of the expected number of parsimony steps
// implied by a bipartition under the hypothesis that the tree is random.
// In Input:
// - The max depth
// - The number of taxa
// - A pointer to a 2D array (given by precompute_steps_probability(int max_depth, int nb_tax)):
//    * First dimension : depth
//    * Second dimension : steps
//    * value : probability of the step at a given depth
// In output : An array with :
// - the depth in index
// - the expected Number of random parsimony steps
func expected_pars_rand_steps(max_depth, nb_tax int, distribution_rand_steps [][]float64) []float64 {
	expected_rand_steps := make([]float64, max_depth+1)

	for depth := 0; depth <= max_depth; depth++ {
		probas := make([]float64, depth+1)
		for cur_step := 0; cur_step <= depth; cur_step++ {
			probas[cur_step] = float64(cur_step) * distribution_rand_steps[depth][cur_step]
		}
		sort.Float64s(probas)
		sum_proba := 0.0
		for i := 0; i <= depth; i++ {
			sum_proba += probas[i]
		}
		expected_rand_steps[depth] = sum_proba - 1
	}
	return expected_rand_steps
}

//   Precomputes the distribution of the number of steps at all depths
//   Output:
//   - A 2D array:
//      * First dimension : depth
//      * Second dimension : steps
//      * value : probability of the step at a given depth
func precompute_pars_steps_probability(max_depth, nb_tax int) [][]float64 {
	distribution_rand_steps := make([][]float64, max_depth+1)

	for depth := 0; depth <= max_depth; depth++ {
		distribution_rand_steps[depth] = make([]float64, max_depth+1)
		for cur_step := 0; cur_step <= depth; cur_step++ {
			distribution_rand_steps[depth][cur_step] = compute_probability_step_depth(cur_step, nb_tax, depth)
		}
	}

	return distribution_rand_steps
}

// This function computes the probability of a given number parsimony steps
// at a given depth in the tree
func compute_probability_step_depth(nb_steps, nb_tax, depth int) float64 {
	output := float64(nb_steps)*math.Log(2) +
		math.Log(2*float64(nb_tax)-3*float64(nb_steps)) +
		gostats.Factorial_log_rmnj(2*depth-nb_steps-1) +
		gostats.Factorial_log_rmnj(2*(nb_tax-depth)-nb_steps-1) +
		gostats.Factorial_log_rmnj(nb_tax-nb_steps) -
		gostats.Factorial_log_rmnj(depth-nb_steps) -
		gostats.Factorial_log_rmnj(nb_tax-depth-nb_steps) -
		gostats.Factorial_log_rmnj(nb_steps-1) -
		gostats.Factorial_log_rmnj(2*nb_tax-2*nb_steps)
	return math.Exp(output)
}

func nbParsimonySteps(e *tree.Edge, bootTree *tree.Tree) int {
	state := STATENULL

	steps := nbParsimonyStepsRecur(bootTree.Root(), nil, bootTree, e, &state, 0)

	return (steps - 1)
}

func maxDepth(edges []*tree.Edge) int {
	max_depth := 0
	for _, e := range edges {
		var d int
		var err error
		if d, err = e.TopoDepth(); err != nil {
			panic(err)
		}
		max_depth = max(d, max_depth)
	}
	return max_depth
}

func nbParsimonyStepsRecur(cur *tree.Node, prev *tree.Node, bootTree *tree.Tree, e *tree.Edge, state *int, level int) int {
	/* does the post order traversal on current Node and its "descendants" (i.e. not including origin, who is a neighbour of current */
	steps := 0
	sum01 := 0
	sum0 := 0
	sum1 := 0

	/* If it is 01, 0, or 1 */
	nextState := STATENULL

	for _, child := range cur.Neigh() {
		if child != prev {
			// If it is a tip node:
			// STATE 0 or 1
			if child.Tip() {
				bitsetindex, err := bootTree.TipIndex(child.Name())
				if err != nil {
					panic(err)
				}
				if e.TipPresent(bitsetindex) {
					nextState = STATE1
				} else {
					nextState = STATE0
				}
			} else {
				steps += nbParsimonyStepsRecur(child, cur, bootTree, e, &nextState, level+1)
			}

			switch nextState {
			case STATE0:
				sum0++
			case STATE1:
				sum1++
			case STATE01:
				sum01++
			}
		}
	}

	if sum0 == sum1 {
		*state = STATE01
	} else if sum1 > sum0 {
		*state = STATE1
	} else {
		*state = STATE0
	}

	steps += min(sum0, sum1)
	return steps
}

// Thread that takes bootstrap trees from the channel,
// computes the number of pars steps for each edges of the ref tree
// and send it to the result channel
func (supporter *parsimonySupporter) ComputeValue(refTree *tree.Tree, empiricalTrees []*tree.Tree, cpu int, empirical bool, edges []*tree.Edge, randEdges [][]*tree.Edge,
	wg *sync.WaitGroup, bootTreeChannel <-chan utils.Trees, stepsChan chan<- bootval, randStepsChan chan<- bootval) {
	func(cpu int) {
		for treeV := range bootTreeChannel {
			for i, e := range edges {
				if e.Right().Tip() {
					continue
				}
				nbsteps := nbParsimonySteps(e, treeV.Tree)
				stepsChan <- bootval{
					nbsteps,
					i,
					false,
				}
				// We compute the empirical steps
				if empirical {
					for j, _ := range empiricalTrees {
						e2 := randEdges[j][i]
						nbstepsrand := nbParsimonySteps(e2, treeV.Tree)
						randStepsChan <- bootval{
							nbstepsrand,
							i,
							nbsteps >= nbstepsrand,
						}
					}
				}
			}
		}
		wg.Done()
	}(cpu)
}

func (supporter *parsimonySupporter) ExpectedRandValues(maxdepth int, nbtips int) []float64 {
	if supporter.expected_rand_steps == nil {
		supporter.distribution_rand_step_val = precompute_pars_steps_probability(maxdepth, nbtips)
		supporter.expected_rand_steps = expected_pars_rand_steps(maxdepth, nbtips, supporter.distribution_rand_step_val)
	}
	return supporter.expected_rand_steps
}

func Parsimony(reftreefile, boottreefile string, empirical bool, cpus int) *tree.Tree {
	var supporter *parsimonySupporter = &parsimonySupporter{}
	return ComputeSupport(reftreefile, boottreefile, empirical, cpus, supporter)
}

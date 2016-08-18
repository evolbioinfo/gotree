package support

import (
	"fmt"
	"github.com/fredericlemoine/gostats"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"math"
	"runtime"
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

type parsteps struct {
	nbsteps int
	edgeid  int
	randsup bool
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
func expected_rand_steps(max_depth, nb_tax int, distribution_rand_steps [][]float64) []float64 {
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
		expected_rand_steps[depth] = sum_proba
	}
	return expected_rand_steps
}

//   Precomputes the distribution of the number of steps at all depths
//   Output:
//   - A 2D array:
//      * First dimension : depth
//      * Second dimension : steps
//      * value : probability of the step at a given depth
func precompute_steps_probability(max_depth, nb_tax int) [][]float64 {
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

func Parsimony(reftreefile, boottreefile string, empirical bool, cpus int) *tree.Tree {
	var reftree *tree.Tree
	var err error
	nbEmpiricalTrees := 20
	maxcpus := runtime.NumCPU()
	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadRefTree(reftreefile); err != nil {
		panic(err)
	}

	// We generate nbEmpirical shuffled trees
	randEdges := make([][]*tree.Edge, nbEmpiricalTrees)
	if empirical {
		for i := 0; i < nbEmpiricalTrees; i++ {
			var randT *tree.Tree
			if randT, err = utils.ReadRefTree(reftreefile); err != nil {
				panic(err)
			}
			randT.ShuffleTips()
			redges := randT.Edges()
			randEdges[i] = redges
		}
	}

	// We read bootstrap trees
	if boottreefile == "none" {
		panic("You must provide a file containing bootstrap trees")
	}
	var nbtrees int
	compareChannel := make(chan utils.Trees, 100)
	go func() {
		if nbtrees, err = utils.ReadCompTrees(boottreefile, compareChannel); err != nil {
			panic(err)
		}
	}()

	tips := reftree.Tips()
	edges := reftree.Edges()

	// Compute max depth
	max_depth := 0
	for _, e := range edges {
		d, err := e.TopoDepth()
		if err != nil {
			panic(err)
		}
		if d > max_depth {
			max_depth = d
		}
	}

	distribution_rand_step_val := precompute_steps_probability(max_depth, len(tips))
	expected_rand_step_val := expected_rand_steps(max_depth, len(tips), distribution_rand_step_val)

	stepsChan := make(chan parsteps, 1000)
	randStepsChan := make(chan parsteps, 1000)

	stepsBoot := make([]int, len(edges))
	var wg sync.WaitGroup
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go func(cpu int) {
			for treeV := range compareChannel {
				for i, e := range edges {
					if e.Right().Tip() {
						continue
					}
					nbsteps := nbParsimonySteps(e, treeV.Tree)
					stepsChan <- parsteps{
						nbsteps,
						i,
						false,
					}
					// We compute the empirical steps
					if empirical {
						for j := 0; j < nbEmpiricalTrees; j++ {
							e2 := randEdges[j][i]
							nbstepsrand := nbParsimonySteps(e2, treeV.Tree)
							randStepsChan <- parsteps{
								nbstepsrand,
								j,
								nbsteps >= nbstepsrand,
							}
						}
					}

				}
			}
			wg.Done()
		}(cpu)
	}

	go func() {
		wg.Wait()
		close(stepsChan)
		close(randStepsChan)
	}()

	gtRandom := make([]int, len(edges))

	var wg2 sync.WaitGroup
	wg2.Add(2)
	go func() {
		for step := range stepsChan {
			stepsBoot[step.edgeid] += step.nbsteps
			d, err := edges[step.edgeid].TopoDepth()
			if err != nil {
				panic(err)
			}
			if !empirical {
				randstep := expected_rand_step_val[d] - 1
				if float64(step.nbsteps) >= randstep {
					gtRandom[step.edgeid]++
				}
			}
		}
		wg2.Done()
	}()

	stepRand := make([]int, len(edges))
	if empirical {
		go func() {
			for step := range randStepsChan {
				if step.randsup {
					gtRandom[step.edgeid]++
				}
				stepRand[step.edgeid] += step.nbsteps
			}
			wg2.Done()
		}()
	} else {
		wg2.Done()
	}

	wg2.Wait()

	for i, e := range edges {
		if !edges[i].Right().Tip() {
			d, err := e.TopoDepth()
			if err != nil {
				panic(err)
			}
			avg_val := float64(stepsBoot[i]) / float64(nbtrees)
			var pval, avg_rand_val float64
			if empirical {
				avg_rand_val = float64(stepRand[i]) / float64(nbEmpiricalTrees)
				pval = float64(gtRandom[i]) * 1.0 / (float64(nbEmpiricalTrees) * float64(nbtrees))
			} else {
				avg_rand_val = expected_rand_step_val[d] - 1
				pval = float64(gtRandom[i]) * 1.0 / float64(nbtrees)
			}
			support := float64(1) - avg_val/avg_rand_val

			edges[i].SetSupport(support)
			edges[i].Right().SetName(fmt.Sprintf("%.2f/%.4f", support, pval))
		}
	}

	return reftree
}

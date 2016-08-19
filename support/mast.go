package support

func MastLike(reftreefile, boottreefile string, empirical bool, cpus int) *tree.Tree {
	var reftree *tree.Tree             // reference tree
	var err error                      // error output
	var nbEmpiricalTrees int = 10      // number of empirical trees to simulate
	var maxcpus int = runtime.NumCPU() // max number of cpus
	var randEdges [][]*tree.Edge       // Edges of shuffled trees
	var nbtrees int                    // Number of bootstrap trees
	var max_depth int                  // Maximum topo depth of all edges of ref tree
	var tips []*tree.Node              // Tip nodes of the ref tree
	var edges []*tree.Edge             // Edges of the reference tree
	var stepsBoot []int                // Sum of number of steps per edge over boot trees
	var stepRand []int                 // Sum of number of steps per random edges over boot trees
	var gtRandom []int                 // Number of times edges have steps that are >= rand steps

	var wg sync.WaitGroup                // For waiting end of step computation
	var wg2 sync.WaitGroup               // For waiting end of final counting
	var stepsChan chan parsteps          // Channel of number of steps computed for a given edge
	var randStepsChan chan parsteps      // Channel of number of steps computed for a given shuffled edge
	var bootTreeChannel chan utils.Trees // Channel of bootstrap trees

	stepsChan = make(chan parsteps, 1000)
	randStepsChan = make(chan parsteps, 1000)
	bootTreeChannel = make(chan utils.Trees, 100)

	if cpus > maxcpus {
		cpus = maxcpus
	}

	if reftree, err = utils.ReadRefTree(reftreefile); err != nil {
		panic(err)
	}
	tips = reftree.Tips()
	edges = reftree.Edges()
	max_depth = maxDepth(edges)
	stepsBoot = make([]int, len(edges))
	gtRandom = make([]int, len(edges))
	stepRand = make([]int, len(edges))

	// Precomputation of expected number of parsimony steps per depth
	distribution_rand_step_val := precompute_steps_probability(max_depth, len(tips))
	expected_rand_step_val := expected_rand_steps(max_depth, len(tips), distribution_rand_step_val)

	// We generate nbEmpirical shuffled trees and store their edges
	randEdges = make([][]*tree.Edge, nbEmpiricalTrees)
	if empirical {
		for i := 0; i < nbEmpiricalTrees; i++ {
			var randT *tree.Tree
			if randT, err = utils.ReadRefTree(reftreefile); err != nil {
				panic(err)
			}
			randT.ShuffleTips()
			randEdges[i] = randT.Edges()
		}
	}

	// We read bootstrap trees and put them in the channel
	if boottreefile == "none" {
		panic("You must provide a file containing bootstrap trees")
	}
	go func() {
		if nbtrees, err = utils.ReadCompTrees(boottreefile, bootTreeChannel); err != nil {
			panic(err)
		}
	}()

	// We compute parsimony steps for all bootstrap trees
	for cpu := 0; cpu < cpus; cpu++ {
		wg.Add(1)
		go mastDistComputation(nbEmpiricalTrees, cpu, empirical, edges, randEdges, &wg, bootTreeChannel, stepsChan, randStepsChan)
	}

	// Wait for step computation to close output channels
	go func() {
		wg.Wait()
		close(stepsChan)
		close(randStepsChan)
	}()

	// Now count steps from the output channels
	wg2.Add(2)
	go func() {
		for step := range stepsChan {
			stepsBoot[step.edgeid] += step.nbsteps
			d, err := edges[step.edgeid].TopoDepth()
			if err != nil {
				panic(err)
			}
			// If theoretical we must count number >= here
			if !empirical {
				randstep := expected_rand_step_val[d] - 1
				if float64(step.nbsteps) >= randstep {
					gtRandom[step.edgeid]++
				}
			}
		}
		wg2.Done()
	}()

	// If "empirical" we read the randStepsChan
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

	// Finally we compute pvalues and support and write it in the tree
	for i, e := range edges {
		if !edges[i].Right().Tip() {
			d, err := e.TopoDepth()
			if err != nil {
				panic(err)
			}
			avg_val := float64(stepsBoot[i]) / float64(nbtrees)
			var pval, avg_rand_val float64
			if empirical {
				avg_rand_val = float64(stepRand[i]) / (float64(nbEmpiricalTrees) * float64(nbtrees))
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

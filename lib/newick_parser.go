package lib

// This parser has been translated from C to GO
// The original was written by Jean Baka Domelevo Entfellner
import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

func FromNewickFile(file string) (*Tree, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	treeString := string(dat)
	tree, err2 := FromNewickString(treeString)
	if err2 != nil {
		return nil, err2
	}
	return tree, nil
}

func FromNewickString(newick_str string) (*Tree, error) {
	treeString := []rune(newick_str)
	tree := NewTree()
	begin, end := 0, 0 /* to delimitate the string to further process */
	// SYNTACTIC CHECKS on the input string
	i := 0
	for {
		if newick_str[i] != ' ' && newick_str[i] != '\t' && newick_str[i] != '\n' && newick_str[i] != '\r' {
			break
		}
		i++
	}

	// begin: AFTER the very first parenthesis
	if newick_str[i] != '(' {
		return nil, errors.New("Tree doesn't start with an opening parenthesis")
	} else {
		begin = i + 1
	}

	i = len(newick_str) - 1
	for {
		if newick_str[i] != ' ' && newick_str[i] != '\t' && newick_str[i] != '\n' && newick_str[i] != '\r' {
			break
		}
		i -= 1
	}

	if newick_str[i] != ';' {
		return nil, errors.New("Tree doesn't end with a semicolon")
	}

	// end: BEFORE the very last parenthesis, discarding optional name for the root and uncanny branch length for its "father" branch
	for {
		if newick_str[i] == ')' {
			break
		}
		i--
	}

	end = i - 1

	err := parseNewickRune(treeString, tree, nil, begin, end, false)
	if err != nil {
		return nil, err
	}
	update_bootstrap_supports_from_node_names(tree, nil, nil)
	tree.UpdateTipIndex()
	tree.clearBitSetsRecur(nil, nil, uint(len(tree.tipIndex)))
	tree.UpdateBitSet()
	return tree, nil
}

func parseNewickRune(newick_str []rune, tree *Tree, curnode *Node, begin int, end int, has_father bool) error {
	if begin > end {
		return errors.New("Error in parse_substring_into_node: begin > end\n")
	}
	if curnode == nil {
		curnode = tree.AddNewNode()
		tree.SetRoot(curnode)
	}
	nb_commas := countOuterCommas(newick_str, begin, end)
	var pair0, pair1, inner_pair0, inner_pair1 int
	comma_index := begin - 1
	if nb_commas > 0 {
		// at least one comma, so at least two sons:
		for i := 0; i <= nb_commas; i++ { /* e.g. three iterations for two commas */
			pair0 = comma_index + 1 /* == begin at first iteration */
			if i == nb_commas {
				comma_index = end + 1
			} else {
				comma_index = index_next_toplevel_comma(newick_str, pair0, end)
			}
			pair1 = comma_index - 1

			son, err := create_son_and_connect_to_father(curnode, tree, newick_str, pair0, pair1)
			if err != nil {
				return err
			}

			// RECURSIVE TREATMENT OF THE SON
			//because name and brlen already processed by create_son
			err = strip_toplevel_parentheses(newick_str, pair0, pair1, &inner_pair0, &inner_pair1)
			if err != nil {
				return err
			}
			// recursive treatment

			err = parseNewickRune(newick_str, tree, son, inner_pair0, inner_pair1, true)
			if err != nil {
				return err
			}
			// after the recursive treatment of the son, the data structures of the son have been created, so now we can write
			//   in it the data corresponding to its direction0 (father) */
			// son->neigh[0] = current_node;
			// son->br[0] = current_node->br[direction];
		} /* end for i (treatment of the various sons) */
	} /* end if/else on the number of commas */
	return nil
}

func countOuterCommas(str []rune, begin int, end int) int {
	/* returns the number of toplevel commas found, from position begin included, up to position end. */
	count, level := 0, 0
	for i := begin; i <= end; i++ {
		switch str[i] {
		case '(':
			level++
		case ')':
			level--
		case ',':
			if level == 0 {
				count++
			}
		} /* endswitch */
	} /* endfor */
	return count
} /* end count_outer_commas */

func strip_toplevel_parentheses(in_str []rune, begin int, end int, pair0 *int, pair1 *int) error {
	// returns the new (begin,end) pair comprising all chars found strictly inside the toplevel parentheses.
	// The input "pair" is an array of two integers, we are passing the output values through it.
	// It is intended that here, in_str[pair[0]-1] == '(' and in_str[pair[1]+1] == ')'.
	// In case no matching parentheses are simply return begin and end in pair[0] and pair[1]. It is NOT an error.
	// This function also tests the correctness of the NH syntax: if no balanced pars, then return an error and abort
	found_par := 0

	*pair0 = end + 1
	*pair1 = -1 /* to ensure termination if no parentheses are found */

	/* first seach opening par from the beginning of the string */
	for i := begin; i <= end; i++ {
		if in_str[i] == '(' {
			*pair0 = i + 1
			found_par += 1
			break
		}
	}

	/* and then search the closing par from the end of the string */
	for i := end; i >= begin; i-- {
		if in_str[i] == ')' {
			*pair1 = i - 1
			found_par += 1
			break
		}
	}

	switch found_par {
	case 0:
		*pair0 = begin
		*pair1 = end
		break
	case 1:
		return errors.New("Syntax error in NH tree: unbalanced parentheses between string indices " + strconv.Itoa(begin) + " and " + strconv.Itoa(end) + ". Aborting.")
	}
	return nil
	// end of switch: nothing to do in case 2 (as pair[0] and pair[1] correctly set), and found_par can never be > 2
}

// This function creates (allocates) the son node in the given direction from the current node.
// It also creates a new branch to connect the son to the father.
// The array structures in the tree (a_nodes and a_edges) are updated accordingly.
// Branch length and node name are processed.
// The input string given between the begin and end indices (included) is of the type:
// (...)node_name:length
// OR
// leaf_name:length
// OR
// a:1,b:0.31,c:1.03
// In both cases the length is optional, and replaced by MIN_BR_LENGTH if absent. */
func create_son_and_connect_to_father(current_node *Node, current_tree *Tree, in_str []rune, begin int, end int) (*Node, error) {

	son := current_tree.AddNewNode()
	edge := current_tree.ConnectNodes(current_node, son)

	err := process_name_and_brlen(son, edge, current_tree, in_str, begin, end)
	if err != nil {
		return nil, err
	}
	return son, nil
}

// looks into in_str[begin..end] for the branch length of the "father" edge
//	   and updates the edge and node structures accordingly */
func process_name_and_brlen(son_node *Node, edge *Edge, current_tree *Tree, in_str []rune, begin int, end int) error {
	colon := index_toplevel_colon(in_str, begin, end)
	closing_par, opening_bracket := -1, -1
	var name_begin, name_end, name_length, effective_length int

	/* processing the optional BRANCH LENGTH... */
	if colon != -1 {
		brlen, err := strconv.ParseFloat(string(in_str[colon+1:end]), 64)
		if err != nil {
			return err
		}
		edge.length = brlen
	}

	/* then scan backwards from the colon (or from the end if no branch length) to get the NODE NAME,
	   not going further than the first closing par */
	/* we ignore the NHX-style comments for the moment, hence the detection of the brackets, which can contain anything but nested brackets */
	var ignore_mode bool = false
	var endpos int
	if colon == -1 {
		endpos = end
	} else {
		endpos = colon - 1
	}
	for i := endpos; i >= begin; i-- {
		if in_str[i] == ']' && !ignore_mode {
			ignore_mode = true
		} else if in_str[i] == ')' && !ignore_mode {
			closing_par = i
			break
		} else if in_str[i] == '[' && ignore_mode {
			ignore_mode = false
			opening_bracket = i
		}
	} /* endfor */

	if closing_par == -1 {
		name_begin = begin
	} else {
		name_begin = closing_par + 1
	}

	if opening_bracket != -1 {
		name_end = opening_bracket - 1
	} else {
		if colon == -1 {
			name_end = end
		} else {
			name_end = colon - 1
		}
	}

	/* but now if the name starts and ends with single or double quotes, remove them */
	if in_str[name_begin] == in_str[name_end] && (in_str[name_begin] == '"' || in_str[name_begin] == '\'') {
		name_begin++
		name_end--
	}
	name_length = name_end - name_begin + 1
	effective_length = name_length
	if name_length >= 1 {
		son_node.name = string(in_str[name_begin : name_begin+effective_length])
	}
	return nil
}

// returns the index of the (first) toplevel colon only, -1 if not found
func index_toplevel_colon(in_str []rune, begin int, end int) int {
	level := 0
	/* more efficient to proceed from the end in this case */
	for i := end; i >= begin; i-- {
		switch in_str[i] {
		case ')':
			level++
		case '(':
			level--
		case ':':
			if level == 0 {
				return i
			}
		}
	}
	return -1
}

// returns the index of the next toplevel comma, from position begin included, up to position end.
//   the result is -1 if none is found.
func index_next_toplevel_comma(in_str []rune, begin int, end int) int {
	level := 0
	for i := begin; i <= end; i++ {
		switch in_str[i] {
		case '(':
			level++
		case ')':
			level--
		case ',':
			if level == 0 {
				return i
			}
		}
	}
	// reached if no outer comma found
	return -1
}

// a branch takes its support value from its descendant node (son).
//	   The current node under examination will give its value (node name) to its father branch, if that one exists.
//	   We modify here the bootstrap support on the edge between current and origin. It is assumed that the node "origin" is on
//	   the path from "current" to the (pseudo-)root */
func update_bootstrap_supports_from_node_names(t *Tree, current *Node, origin *Node) error {
	if current == nil {
		current = t.Root()
	}

	if origin != nil {
		edgei, err := current.NodeIndex(origin)
		if err != nil {
			return err
		}
		edge := current.br[edgei]
		if current.name != "" {
			boot, err := strconv.ParseFloat(current.name, 64)
			if err == nil {
				/* if succesfully parsing a number */
				edge.support = boot
				current.name = ""
				fmt.Println(boot)
			}
		}
	}

	for _, child := range current.neigh {
		if child != origin {
			update_bootstrap_supports_from_node_names(t, child, current)
		}
	}
	return nil
}

package newick_test

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io/newick"
	"strings"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseTree(t *testing.T) {
	intree := "(Tip2:1.00000,(Tip 7:1.00000,((Tip9[COUCOU TOUT LE MONDE]:1.00000,((Tip5:1.00000,((Tip8:1.00000,Tip6:0.50000):0.50000,Tip4:0.25000):0.25000):0.50000,Tip3:0.50000):0.25000):0.25000,Node0:0.25000):0.12500):0.12500,Tip1:0.50000);"
	_, err := newick.NewParser(strings.NewReader(intree)).Parse()
	if err != nil {
		fmt.Println(err)
	}
}

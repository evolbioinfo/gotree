// package mutations provides data structures & functions
// for counting mutations given an alignment of ancestral and tips sequences
package mutations

import (
	"fmt"

	"github.com/evolbioinfo/gotree/io"
)

type Mutation struct {
	AlignmentSite             int    // Index of the site of the alignment
	BranchIndex               int    // Index of the branch
	ChildNodeName             string // Name of the parent of the clade
	ParentCharacter           uint8  // Parent character
	ChildCharacter            uint8  // Child character
	NumTips                   int    // Total number of descendent tips
	NumTipsWithChildCharacter int    // Number of descendent tips that have the child character
	NumEEM                    int    // Number of emergence of this mutation
}

type MutationList struct {
	Mutations map[string]Mutation // Key: "AlignmentSite-BranchIndex-ParentCharacter-ChildCharacter"
}

func NewMutationList() (mutations *MutationList) {
	mutations = &MutationList{
		Mutations: make(map[string]Mutation),
	}
	return
}

func (m *MutationList) Exists(id string) (exist bool) {
	_, exist = m.Mutations[id]
	return exist
}

func (m *MutationList) Append(mapp *MutationList) (err error) {
	var exist bool
	for k, v := range mapp.Mutations {
		if _, exist = m.Mutations[k]; exist {
			err = fmt.Errorf("mutation %s already exist in the list", k)
			io.LogError(err)
			return
		}
		m.Mutations[k] = v
	}
	return
}

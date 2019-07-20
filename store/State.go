package store

import (
	"github.com/tendermint/iavl"
)

const (
	accountPrefix = "account::"
	secretPrefix  = "secret::"
	medicalRecord = "medicalrecord::"
)

//State ...
type State struct {
	Tree iavl.MutableTree
}

//NewStateFromTree ...
func NewStateFromTree(tree iavl.MutableTree) *State {
	return &State{tree}
}
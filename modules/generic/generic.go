package generic

import (
	"github.com/astr0n8t/gopen/definitions"
)

type Generic struct {
	Result      bool
	Output      string
	Options     definitions.Options
	Process     definitions.Process
	ResultStore definitions.ResultStore
}

func New(opts definitions.Options, proc definitions.Process, res definitions.ResultStore) *Generic {
	return &Generic{false, "", opts, proc, res}
}

func (n *Generic) RunModule() definitions.ResultStore {
	n.Result = true
	return n.ResultStore
}

func (n *Generic) GetOutput() string {
	return n.Output
}

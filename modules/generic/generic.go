package generic

import (
	"github.com/astr0n8t/gopen/definitions"
)

type Generic struct {
	Result bool
	Output string
}

func New(opts definitions.Options, proc definitions.Process) *Generic {
	return &Generic{false, ""}
}

func (n *Generic) RunModule() bool {
	n.Result = true
	return n.Result
}

func (n *Generic) GetOutput() string {
	return n.Output
}

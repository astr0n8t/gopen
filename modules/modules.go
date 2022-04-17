package modules

import (
	"github.com/astr0n8t/gopen/definitions"
	"github.com/astr0n8t/gopen/modules/generic"
	"github.com/astr0n8t/gopen/modules/nmap"
)

// Module interface defines what functions a provider struct should have
type Module interface {
	RunModule() definitions.ResultStore
	GetOutput() string
}

// GetModule returns a new object of the given module type
func GetModule(name string, options definitions.Options, process definitions.Process, result definitions.ResultStore) Module {

	if name == "nmap" {
		return nmap.New(options, process, result)
	}

	return generic.New(options, process, result)
}

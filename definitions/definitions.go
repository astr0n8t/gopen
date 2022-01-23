package definitions

// A struct to store configuration options
type Config struct {
	Variables Options
	Workflow  map[string]Process
}

type Options struct {
	Addresses string
	Ports     string
	Root      bool
}

type Process struct {
	Executable string
	Flags      string
}

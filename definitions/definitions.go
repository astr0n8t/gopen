package definitions

// A struct to store configuration options
type Config struct {
	Variables Options
	Workflow  map[string]Process
}

type Options struct {
	Threads   int
	Addresses string
	Ports     string
	Root      bool
}

type Process struct {
	Executable string
	Flags      string
}

type ResultStore struct {
	Hosts    map[string]Host
	Previous string
	Success  bool
}

type Host struct {
	Address     string
	Port        map[int]bool
	OS          string
	Fingerprint string
	Pages       []string
	Misc        string
}

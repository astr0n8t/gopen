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
	Hostnames   []string
	Ports       map[int]Port
	OS          string
	Fingerprint string
	Pages       []string
	Misc        string
}

type Port struct {
	ID       int
	Protocol string
	State    string
}

type ThreadResult struct {
	Result bool
	Err    string
}

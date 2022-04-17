package nmap

import "github.com/astr0n8t/gopen/definitions"

type Nmap struct {
	Result      bool
	Output      string
	Options     definitions.Options
	Process     definitions.Process
	ResultStore definitions.ResultStore
}

func New(opts definitions.Options, proc definitions.Process, res definitions.ResultStore) *Nmap {
	// Return the reference to a new Cloudflare object
	return &Nmap{false, "", opts, proc, res}
}

func (n *Nmap) RunModule() definitions.ResultStore {
	n.Result = true
	return n.ResultStore
}

func (n *Nmap) GetOutput() string {
	return n.Output
}

/* Previous POC code
func nmap(send chan string, address string, port int, proc definitions.Process) {
	cmd := exec.Command("nmap", address, "-p", strconv.Itoa(port))
	out, _ := cmd.Output()
	send <- string(out[:])
}


startPort, _ := strconv.Atoi(confOptions.Variables.Ports)

output := make(chan string)

for i := 0; i < 1001; i++ {
	go nmap(output, confOptions.Variables.Addresses, startPort, confOptions.Workflow["nmap"])
	startPort++
}

msg := <-output
for i := 0; i < 1000; i++ {
	msg += <-output
}

fmt.Println(msg)
*/

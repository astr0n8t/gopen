package nmap

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/astr0n8t/gopen/definitions"
)

type Nmap struct {
	Result      bool
	Output      string
	Options     definitions.Options
	Process     definitions.Process
	ResultStore *definitions.ResultStore
}

// nmapXML allows unmarshalling the output of NMAP into something useful
// note: we should only ever process one host at a time
// needs updated for other things that can show up with a -A
type nmapXML struct {
	XMLName xml.Name `xml:"nmaprun"`
	Host    struct {
		Hostnames []struct {
			Hostname struct {
				Name string `xml:"name,attr"`
			} `xml:"hostname"`
		} `xml:"hostnames"`
		Ports []struct {
			Port struct {
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					State string `xml:"state,attr"`
				} `xml:"state"`
				Service struct {
					Name string `xml:"name,attr"`
				} `xml:"service"`
			} `xml:"port"`
		} `xml:"ports"`
	} `xml:"host"`
}

func New(opts definitions.Options, proc definitions.Process, res *definitions.ResultStore) *Nmap {
	// Return the reference to a new nmap object
	return &Nmap{false, "", opts, proc, res}
}

// Returns the pointer to the updated ResultStore object
// Resets the module result
// Calls the spawn function
// Public Function
func (n *Nmap) RunModule() *definitions.ResultStore {
	n.Result = true
	n.spawn()
	n.ResultStore.Success = n.Result
	return n.ResultStore
}

// stub
// Returns the textual update of the module
// Public function
func (n *Nmap) GetOutput() string {
	return n.Output
}

// Handles the spawning of the nmap threads
// needs to be updated to take ports into account for hosts of smaller size
// Private function
// Updates Result with any errors and ResultStore with any discovered information
func (n *Nmap) spawn() {
	startPort, _ := strconv.Atoi(n.Options.Ports)

	errorChannels := make(chan definitions.ThreadResult)
	output := make(chan definitions.Host)

	iterations := len(n.ResultStore.Hosts)/n.Options.Threads + 1
	remainder := len(n.ResultStore.Hosts) % n.Options.Threads

	hostKeys := make([]string, len(n.ResultStore.Hosts))

	// Get an array of the keys so we can keep track of what we have done
	currentKey := 0
	for host := range n.ResultStore.Hosts {
		hostKeys[currentKey] = host
		currentKey++
	}

	spawnNum := n.Options.Threads

	// Enter the main loop
	for i := 0; i < iterations; i++ {
		// Check if we need to spawn a full set of threads or just the remaining threads
		if i == iterations-1 {
			spawnNum = remainder
		}
		// Spawn the necessary number of threads
		for j := 0; j < spawnNum; j++ {
			go nmap(output, errorChannels, n.ResultStore.Hosts[hostKeys[i+j]], startPort, n.Process)
		}

		// Receive the output from all of the threads
		for j := 0; j < spawnNum; j++ {
			// Retrieve the updated host object and update results
			tmpHost := <-output
			n.ResultStore.Hosts[tmpHost.Address] = tmpHost

			// Perform error handling
			tmpResult := <-errorChannels
			if !tmpResult.Result {
				fmt.Printf("Error in nmap thread #%d: %s\n", i+j, tmpResult.Err)
				// Update module result
				n.Result = false
			}
		}
	}
}

// The function to run nmap concurrently
// Will try to run nmap and store the results in the options.Host struct datatype
// Input:
// send - Channel for sending the host
// issue - Channel for error handling
// host - The host to check
// port - the port to check
// proc - the options nmap should be run with
// Output:
// Sends back a populated host object via the send channel
// Can also report back errors from a thread via the issue channel
func nmap(send chan definitions.Host, issue chan definitions.ThreadResult, host definitions.Host, port int, proc definitions.Process) {
	// Call nmap
	cmd := exec.Command(proc.Executable, proc.Flags, "-p", strconv.Itoa(port), host.Address, "-oX", "-")
	// Get output and errors
	out, err := cmd.CombinedOutput()

	// If there are errors, send the errors and exit
	if err != nil {
		send <- host
		issue <- definitions.ThreadResult{
			Result: false,
			Err:    string(out)}
		return
	}

	// Unmarshall the output as xml
	var xmlOutput nmapXML
	xml.Unmarshal(out, &xmlOutput)

	// Add the hostnames to the host object
	for _, hostname := range xmlOutput.Host.Hostnames {
		host.Hostnames = append(host.Hostnames, hostname.Hostname.Name)
	}

	// Add the ports scanned to the host object
	for _, port := range xmlOutput.Host.Ports {
		currentPort := definitions.Port{}
		currentPort.ID, _ = strconv.Atoi(port.Port.Portid)
		currentPort.Protocol = port.Port.Protocol
		currentPort.State = port.Port.State.State
		host.Ports[currentPort.ID] = currentPort
	}

	// Send back the updated host object
	// Send back that no errors occurred
	send <- host
	issue <- definitions.ThreadResult{Result: true}
}

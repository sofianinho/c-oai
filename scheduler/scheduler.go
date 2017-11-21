package scheduler

import(
	"syscall"
	"github.com/sofianinho/vnf-api-golang/vnf/types"
)

//State type represents the tasks possible state for the scheduler
// created, running, exited, and stopped are the possible states of the task
type State uint8
const(
	created	State = iota
	running
	exited
	stopped
)
var states = []string{"created", "running", "exited", "stopped"}
// String() returns the state as a string
func (s State) String() string {
	return states[s]
}

//API is the interface to interact with your schedulers. Supported ones are Linux (system) and Docker (docker). You can implement your own simply around the same interface.
type API interface{
	InstanceRun(*types.Instance)(error)
	InstanceKill(sig syscall.Signal, iID string)(error)
	InstanceStatus(iID string)(string,error)

	SystemStatus()(*types.Status,error)
}
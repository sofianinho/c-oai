package scheduler

import (
	"sync"
	"github.com/sofianinho/vnf-api-golang/vnf/types"
)

type taskDocker struct{
	instance	*types.Instance
	status		State
	containerID	string
}

//Containers represents the list of instances during runtime on the docker daemon. 
type Containers struct{
	rw		sync.Mutex
	list	map[string]*taskDocker
}


package scheduler

import (
	"path/filepath"
	"fmt"
	"sync"
	"syscall"
	"os/exec"
	"os"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/utils"
)

type taskSystem struct{
	instance	*types.Instance
	status		State
	pid			int
}


//New returns a system (linux) backed scheduler interface
func New()(API, error){
	return &Tasks{list: map[string]*taskSystem{}}, nil
}

//Tasks represents the list of instances during runtime. 
//The ID of the tasks map is the ID of the instance.
type Tasks struct{
	rw		sync.Mutex
	list	map[string]*taskSystem
}

//InstanceRun actually running the instance on the system
func (t *Tasks)InstanceRun(i *types.Instance)(error){
	t.rw.Lock()
	defer t.rw.Unlock()
	_, ok := t.list[i.ID]
	//this is the first time the instance is going to run
	if !ok{
		k := &taskSystem{i, created, 0}
		t.list[i.ID] = k
	}else{
		//the task already exists in the list
		t.list[i.ID].status = running
	}
	//TODO create the task to run using exec.Command and exec.Start and update the pid/status
	//1. Check path
	if _, err := exec.LookPath(i.Artefact); err != nil{
		fmt.Println("Path is: ", os.Getenv("PATH"))
		return fmt.Errorf("Task %s has artefact %s which is not in PATH: %s", i.ID, i.Artefact,err)
	}
	//2. Run task
	lF := filepath.Join(config.Params.GetString("runtime.path"), i.ID, i.Artefact+".log")
	logFile, e := os.Create(lF)
	if e != nil{
		return fmt.Errorf("failed to open log file %s: %s", lF, e)
	}
	defer logFile.Close()
	lEF := filepath.Join(config.Params.GetString("runtime.path"), i.ID, i.Artefact+".err")
	logErrFile, err := os.Create(lEF)
	if err != nil{
		return fmt.Errorf("failed to open error file %s: %s", lEF, err)
	}
	defer logErrFile.Close()
	wd, e := os.Getwd()
	if e != nil{
		return fmt.Errorf("Unable to get current dir: %s", e)
	}
	confFile := filepath.Join(wd, config.Params.GetString("runtime.path"), i.ID, "enb.conf")
	fmt.Println("conf file: ", confFile)
	cmd := exec.Command(i.Artefact, "--ulsch-max-errors", "100000", "-S", "-O",confFile)
	//3. Capture stdout and stderr of command
	cmd.Stdout = logFile
	cmd.Stderr = logErrFile
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("could not run task %s: %s", i.ID, err)
	}
	
	t.list[i.ID].pid = cmd.Process.Pid
	t.list[i.ID].status = running
	return nil
}

//InstanceKill send a signal to a running instance on the system
func (t *Tasks)InstanceKill(sig syscall.Signal, iID string)(error){
	t.rw.Lock()
	defer t.rw.Unlock()

	k, ok := t.list[iID]
	if !ok{
		return fmt.Errorf("Instance %s not in the tasks list", iID)
	}
	if k.status == running{
		if err := syscall.Kill(k.pid, sig); err!=nil{
			return fmt.Errorf("Error while sending signal to instance %s: %s",iID, err)
		}
		if (sig == syscall.SIGKILL)||(sig == syscall.SIGINT) {
			syscall.Kill(t.list[iID].pid, sig)
		}
	}
	k.status = stopped
	return nil
}

//InstanceStatus returns the running status of the instance on the system
func (t *Tasks)InstanceStatus(iID string)(string,error){
	t.rw.Lock()
	defer t.rw.Unlock()
	k, ok := t.list[iID]
	if !ok{
		return "", fmt.Errorf("Instance %s not in the tasks list", iID)
	}
	return k.status.String(), nil
}

//SystemStatus returns the global status of the system (load,cpu,ram,net)
func (t *Tasks)SystemStatus()(*types.Status,error){
	return utils.SysStatus()
}

package vnf

import(
	"time"
	"syscall"
	"fmt"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/twinj/uuid"
)

//SessionNew creates a new session on the storage and the system
func (v *vnf)SessionNew()(*types.Session, error){
	defer observeAction("SessionNew", time.Now())
	s:=&types.Session{}
	s.ID = uuid.NewV4().String()
	s.CreatedAt = time.Now()
	s.Instances = map[string]*types.Instance{}
	s.Configs = map[string]*types.Config{}
	config.Log.Infof("New session %s created", s.ID)
	if err := v.storage.SessionSave(s); err!=nil{
		config.Log.Errorf("Error new session: %s", err)
		return nil, err
	}
	v.setGauges()
	return s,nil
}

//SessionGet returns a session given its ID
func (v *vnf)SessionGet(id string)(*types.Session){
	defer observeAction("SessionGet", time.Now())
	if s, e := v.storage.SessionGet(id); e == nil{
		return s
	}
	return nil
}

//SessionDelete deletes a session given its ID
func (v *vnf)SessionDelete(id string)(error){
	defer observeAction("SessionDelete", time.Now())
	
	s := v.SessionGet(id)
	if s == nil{
		return fmt.Errorf("Session %s not found", id)
	}

	s.Lock()
	defer s.Unlock()

	config.Log.Infof("Deleting session %s from the session", id)
	//delete the configs from the session then the instances 
	for _,c := range s.Configs{
		err := v.storage.ConfigDelete(id, c.ID)
		if err != nil{
			config.Log.Errorf("Deleting config %s from session %s failed: %s", c.ID, id, err)
		}
	}
	for _, i := range s.Instances{
		//kill instances from runtime
		err := v.instances.InstanceKill(syscall.SIGINT, i.ID)
		if err != nil{
			config.Log.Errorf("[scheduler] Killing instance %s from session %s failed: %s", i.ID, id, err)
		}
		//delete instances from storage
		err = v.storage.InstanceDelete(id, i.ID)
		if err != nil{
			config.Log.Errorf("[storage] Deleting instance %s from session %s failed: %s", i.ID, id, err)
		}
	}
	return nil
}

//SessionCount returns the number of current sessions
func (v *vnf)SessionCount()(int,error){
	defer observeAction("SessionCount", time.Now())
	return v.storage.SessionCount()
}

//SessionList returns the current sessions
func (v *vnf)SessionList()([]*types.Session, error){
	defer observeAction("SessionList", time.Now())
	return v.storage.SessionList()
}

//SessionStatus returns the current sessions' status
func (v *vnf)SessionStatus()(*types.Status, error){
	defer observeAction("SessionStatus", time.Now())
	return v.instances.SystemStatus()
}
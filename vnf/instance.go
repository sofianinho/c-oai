package vnf

import (
	"syscall"
	"fmt"
	"time"
	"crypto/rand"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
)

const instanceIDLen = 12

//InstanceNew add the instance to storage and then instantiates it on the scheduler
func (v *vnf)InstanceNew(sID string, conf *types.Config, alias string, tags []string, artefact string)(*types.Instance, error){
	defer observeAction("InstanceNew", time.Now())

	i := &types.Instance{}
	key := make([]byte, instanceIDLen)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("Error generating instance ID: %s", err)
	}
	i.ID = fmt.Sprintf("%x", key)
	i.Alias = alias
	i.Session = sID
	i.CreatedAt = time.Now()
	i.Tags = tags
	i.Artefact = artefact
	i.Config = conf
	//save it on disk
	if err := v.storage.InstanceSave(sID,i); err != nil{
		return nil, fmt.Errorf("Could not save new instance in storage: %s", err)
	}
	//update the conf with the ID of the instance that depends on it
	/* if err := v.storage.ConfigAddInstance(sID, conf.ID, i); err!=nil{
		return nil, fmt.Errorf("Could not save new instance to config in storage: %s", err)
	} */

	return i,nil
}

//InstanceGet returns an instance given its ID (if existing)
func (v *vnf)InstanceGet(sID, iID string)(*types.Instance, error){
	defer observeAction("InstanceGet", time.Now())
	return v.storage.InstanceGet(sID, iID)
}

//InstanceGetByAlias returns an instance given its alias (if existing)
func (v *vnf)InstanceGetByAlias(sID,alias string)(*types.Instance, error){
	defer observeAction("InstanceGetByAlias", time.Now())
	return v.storage.InstanceFindByAlias(sID, alias)
}

//InstanceStatus returns the runtime status of an instance(if existing)
func (v *vnf)InstanceStatus(sID,iID string)(string, error){
	defer observeAction("InstanceStatus", time.Now())
	return v.instances.InstanceStatus(iID)
}

//InstanceUpdate updates the instance on disk, kills the runtime instance if config has changed and starts a new with new config. 
//It also updates the previous configuration that pointed to this instance
func (v *vnf)InstanceUpdate(sID,iID string, conf *types.Config, alias string, tags []string, artefact string)(error){
	defer observeAction("InstanceUpdate", time.Now())
	//Update the instance on disk
	i, err := v.storage.InstanceGet(sID, iID)
	if err != nil{
		return fmt.Errorf("Could not update the instance: %s", err)
	}

	//check for colliding aliases in the same session
	t,_ := v.storage.InstanceFindByAlias(sID, alias)
	if t != nil && t.ID != iID{
		return fmt.Errorf("Could not update instance %s with alias %s. Instance %s already has it", iID, alias, t.ID)
	}
	//update the instance
	i.Alias = alias
	i.Tags = tags
	killRunIns := false
	//if the instance's config or artefact changed, update the dependency with old config and kill the running instance to replace with new config or artefact
	//no need to interrupt the running config if config and artefact are the same
	if i.Config.ID != conf.ID || i.Artefact != artefact {
		killRunIns = true
		//update the conf with the ID of the instance that depends on it
		/* if err := v.storage.ConfigDelInstance(sID, conf.ID, i); err!=nil{
			return fmt.Errorf("Could not delete instance from old config in storage: %s", err)
		}
		if err := v.storage.ConfigAddInstance(sID, conf.ID, i); err!=nil{
			return fmt.Errorf("Could not add instance to new config in storage: %s", err)
		} */
		i.Artefact = artefact
		i.Config = conf
	}
	//remove the previous instance from the storage...
	if err := v.storage.InstanceDelete(sID, iID); err != nil{
		return fmt.Errorf("Could not remove old instance from storage: %s", err)
	}
	//...and replace it with the new one
	if err := v.storage.InstanceSave(sID,i); err != nil{
		return fmt.Errorf("Could not save new modifications of instance to storage: %s", err)
	}
	//we have runtime to deal with
	if killRunIns{
		if err := v.instances.InstanceKill(syscall.SIGINT, iID); err!=nil{
			return fmt.Errorf("Could not kill previous running instance %s: %s", iID, err)
		}
		if err := v.instances.InstanceRun(i); err!=nil{
			return fmt.Errorf("Could not start new instance %s: %s", iID, err)
		}
	}
	return nil
}

//InstanceDelete deletes the instance from disk, runtime, and related configs
func (v *vnf)InstanceDelete(sID, iID string)(error){
	defer observeAction("InstanceDelete", time.Now())
	//delete from related configs
	_, err := v.storage.InstanceGet(sID, iID)
	if err != nil{
		return fmt.Errorf("Could not delete the instance: %s", err)
	}
	/* if err := v.storage.ConfigDelInstance(sID, i.Config.ID, i); err!=nil{
		return fmt.Errorf("Could not delete instance from config in storage: %s", err)
	} */
	//delete from runtime
	if err := v.instances.InstanceKill(syscall.SIGINT, iID); err!=nil{
		return fmt.Errorf("Could not kill running instance %s: %s", iID, err)
	}
	//delete from disk
	return v.storage.InstanceDelete(sID, iID)
}

//InstanceCount returns the number of instances in a session
func (v *vnf)InstanceCount(sID string)(int, error){
	defer observeAction("InstanceCount", time.Now())
	return v.storage.InstanceCount(sID)
}

//InstanceList returns the current instances in a session
func (v *vnf)InstanceList(sID string)([]*types.Instance, error){
	defer observeAction("InstanceList", time.Now())
	return v.storage.InstanceList(sID)
}

//InstanceRun runs an instance using the scheduler API
func (v *vnf)InstanceRun(ins *types.Instance)(error){
	defer observeAction("InstanceRun", time.Now())
	return v.instances.InstanceRun(ins)
}

//InstanceRun sends a signal to an instance using the scheduler API
func (v *vnf)InstanceKill(sig syscall.Signal, iID string)(error){
	defer observeAction("InstanceKill", time.Now())
	return v.instances.InstanceKill(sig, iID)
}

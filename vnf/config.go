package vnf

import(
	"fmt"
	"time"
	"crypto/rand"

	"github.com/sofianinho/vnf-api-golang/vnf/types"

)

const configIDLen = 8

//ConfigNew builds a new configuration given its parameters
func (v *vnf)ConfigNew(sID, ver string, conf *types.VNFParams, alias string, tags []string)(*types.Config, error){
	defer observeAction("ConfigNew", time.Now())

	c := &types.Config{}
	key := make([]byte, configIDLen)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("Error generating config ID: %s", err)
	}
	c.ID = fmt.Sprintf("%x", key)
	c.Alias = alias
	c.Session = sID
	c.Version = ver
	//c.Instances = map[string]*types.Instance{}
	c.CreatedAt = time.Now()
	c.Tags = tags
	c.Content = conf
	//save it on disk
	if err := v.storage.ConfigSave(sID,c); err != nil{
		return nil, fmt.Errorf("Could not save new config in storage: %s", err)
	}

	return c,nil
}

//ConfigGet returns a configuration given its id
func (v *vnf)ConfigGet(sID,cID string)(*types.Config, error){
	defer observeAction("ConfigGet", time.Now())

	return  v.storage.ConfigGet(sID, cID)
}

//ConfigGetByAlias returns a configuration given its session ID and alias
func (v *vnf)ConfigGetByAlias(sID,alias string)(*types.Config, error){
	defer observeAction("ConfigGetByAlias", time.Now())
	
	return  v.storage.ConfigFindByAlias(sID,alias)
}

//ConfigUpdate updates the configuration given a session ID and config ID. The updated fields are
// the alias, the tags, and the configuration parameters
func (v *vnf)ConfigUpdate(sID, cID, ver string, conf *types.VNFParams, alias string, tags []string)(error){
	defer observeAction("ConfigUpdate", time.Now())
	
	c, err := v.storage.ConfigGet(sID, cID)
	if err != nil{
		return fmt.Errorf("Could not update the config: %s", err)
	}
	//check for colliding aliases in the same session
	t,_ := v.storage.ConfigFindByAlias(sID, alias)
	if t != nil && t.ID != cID{
		return fmt.Errorf("Could not update config %s with alias %s. Config %s already has it", cID, alias, t.ID)
	}
	//update the config
	c.Alias = alias
	c.Tags = tags
	c.Version = ver
	c.Content = conf
	//remove the previous conf from the storage...
	if err := v.storage.ConfigDelete(sID, cID); err != nil{
		return fmt.Errorf("Could not remove old config from storage: %s", err)
	}
	//...and replace it with the new one
	if err := v.storage.ConfigSave(sID,c); err != nil{
		return fmt.Errorf("Could not save new modifications of config to storage: %s", err)
	}
	//TODO Actually regenerate a configuration file given the config we just updated (templates and all)
	//TODO Restart all the instances that depend on this configuration that changed (especially if some of them are running)
	return nil
}

//ConfigDelete deletes a configuration from a session. 
//This does not affect the instances running on this config. This case should may be handled later.

func (v *vnf)ConfigDelete(sID, cID string)(error){
	defer observeAction("ConfigDelete", time.Now())
	
	//TODO Actually delete the configuration files generated with the config and may be impact running instances
	return v.storage.ConfigDelete(sID, cID)
}

//ConfigCount returns the number of configurations in the session
func (v *vnf)ConfigCount(sID string)(int,error){
	defer observeAction("ConfigCount", time.Now())

	return v.storage.ConfigCount(sID)
}

//ConfigList returns the current configs in a session
func (v *vnf)ConfigList(sID string)([]*types.Config, error){
	defer observeAction("ConfigList", time.Now())
	return v.storage.ConfigList(sID)
}
package storage

import (
	"fmt"
	"sync"
	"encoding/json"
	"os"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/utils"
)

type fileStorage struct{
	rw		sync.Mutex
	path	string
	db		map[string]*types.Session
}


//New creates a new file storage 
func New (p string)(API, error){
	s := &fileStorage{path: p}
	err := s.load()
	if err != nil{
		return nil, err
	}
	
	return s, nil
}

func (s *fileStorage) load() error {
	file, err := os.Open(s.path)
	
	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&s.db)
		if err != nil {
			return err
		}
	} else {
		s.db = map[string]*types.Session{}
	}

	file.Close()
	return nil
}

func (s *fileStorage) save() error {
	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(&s.db)
}

//SessionGet returns the session from file storage
func (s *fileStorage)SessionGet(id string)(*types.Session,error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[id]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", id)
	}
	return sn,nil
}

//SessionSave saves the session in a file storage
func (s *fileStorage)SessionSave(session *types.Session)(error){
	s.rw.Lock()
	defer s.rw.Unlock()

	s.db[session.ID] = session
	return s.save()
}

//SessionCount returns the number of sessions in the storage
func (s *fileStorage)SessionCount()(int,error){
	s.rw.Lock()
	defer s.rw.Unlock()

	return len(s.db),nil
}

//SessionDelete deletes a session given the id
func (s *fileStorage)SessionDelete(id string)(error){
	s.rw.Lock()
	defer s.rw.Unlock()
	delete(s.db, id)
	return nil
}

//SessionList returns all sessions in storage
func (s *fileStorage)SessionList()([]*types.Session, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	var ret []*types.Session

	for _,v := range s.db{
		ret = append(ret, v)
	}

	return ret, nil
}

//ConfigGet returns a config if the sessionID and configID exist
func (s *fileStorage)ConfigGet(sID, cID string)(*types.Config,error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}

	c, ok := sn.Configs[cID]
	if !ok{
		return nil, fmt.Errorf("Config %s in session %s not found",cID,sID)	
	}
	
	return c, nil
}

//ConfigFindByAlias returns a configuration given its alias
func (s *fileStorage)ConfigFindByAlias(sID, alias string)(*types.Config, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}
	for k, v := range sn.Configs{
		if v.Alias == alias{
			return sn.Configs[k], nil
		}
	}
	return nil, fmt.Errorf("Config alias %s not found in session %s",alias,sID)
}

//ConfigFindByTags finds a configuration given its tags
func (s *fileStorage)ConfigFindByTags(sID string, tags []string)([]*types.Config, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	var ret []*types.Config
	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}
	for _, v := range sn.Configs{
		if utils.Includes(v.Tags, tags){
			ret = append(ret, v)
		}
	}
	return ret,nil
}

//ConfigAddInstance adds an instance to a configuration in the storage
/* func (s *fileStorage)ConfigAddInstance(sID, cID string, i *types.Instance)(error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return fmt.Errorf("Session %s not found", sID)
	}

	c, ok := sn.Configs[cID]
	if !ok{
		return fmt.Errorf("Config %s in session %s not found", cID, sID)
	}
	c.Instances[i.ID] = i
	return nil
} */

//ConfigDelInstance deletes an instance from a configuration in the storage
/* func (s *fileStorage)ConfigDelInstance(sID, cID string, i *types.Instance)(error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return fmt.Errorf("Session %s not found", sID)
	}
	c, ok := sn.Configs[cID]
	if !ok{
		return fmt.Errorf("Config %s in session %s not found", cID, sID)
	}

	delete(c.Instances, i.ID)
	return nil
} */


//ConfigSave saves a configuration into a session
func (s *fileStorage)ConfigSave(sID string, c *types.Config)(error){
	s.rw.Lock()
	

	sn, ok := s.db[sID]
	if !ok{
		return fmt.Errorf("Session %s not found", sID)
	}
	sn.Configs[c.ID] = c
	s.rw.Unlock()
	return s.SessionSave(sn)
}

//ConfigCount returns the number of configurations into a session
func (s *fileStorage)ConfigCount(sID string)(int,error){
	s.rw.Lock()
	defer s.rw.Unlock()
	sn, ok := s.db[sID]
	if !ok{
		return 0, fmt.Errorf("Session %s not found", sID)
	}
	return len(sn.Configs),nil
}

//ConfigsCount returns the total number of configs in the storage (all sessions combined)
func (s *fileStorage)ConfigsCount()(int,error){
	s.rw.Lock()
	defer s.rw.Unlock()

	var total int
	for k := range s.db{
		total += len(s.db[k].Configs)
	}
	return total, nil
}

//ConfigDelete deletes a configuration from a session
func (s *fileStorage)ConfigDelete(sID, cID string)(error){
	s.rw.Lock()

	sn, ok := s.db[sID]
	if !ok{
		return fmt.Errorf("Session %s not found", sID)
	}
	delete(sn.Configs, cID)
	s.rw.Unlock()
	return s.SessionSave(sn)
}

//ConfigList returns all configs in a session in storage
func (s *fileStorage)ConfigList(sID string)([]*types.Config, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}
	
	var ret []*types.Config

	for _,v := range sn.Configs{
		ret = append(ret, v)
	}

	return ret, nil
}

//InstanceGet returns an instance state 
func (s *fileStorage)InstanceGet(sID, iID string)(*types.Instance,error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}

	i, ok := sn.Instances[iID]
	if !ok{
		return nil, fmt.Errorf("Instance %s in session %s not found", iID,sID)
	}

	return i,nil
}

//InstanceFindByAlias returns an instance given its alias
func (s *fileStorage)InstanceFindByAlias(sID, alias string)(*types.Instance, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}

	for _, v := range sn.Instances{
		if v.Alias == alias {
			return v, nil
		}
	}

	return nil, fmt.Errorf("Instance alias %s in session %s not found", alias,sID)
}

//InstanceFindByTags returns a slice of instances that share the given tags
func (s *fileStorage)InstanceFindByTags(sID string, tags []string)([]*types.Instance, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	var ret []*types.Instance

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}
	
	for _, v := range sn.Instances{
		if utils.Includes(v.Tags, tags){
			ret = append(ret, v)
		}
	}
	return ret,nil
}

//InstanceSave saves the state of an instance into a session
func (s *fileStorage)InstanceSave(sID string, i *types.Instance)(error){
	s.rw.Lock()

	sn, ok := s.db[sID]
	if !ok{
		return fmt.Errorf("Session %s not found", sID)
	}
	sn.Instances[i.ID] = i
	s.rw.Unlock()
	return s.SessionSave(sn)

}

//InstanceCount returns the number of instances in a session
func (s *fileStorage)InstanceCount(sID string)(int,error){
	s.rw.Lock()
	defer s.rw.Unlock()
	sn, ok := s.db[sID]
	if !ok{
		return 0, fmt.Errorf("Session %s not found", sID)
	}
	return len(sn.Instances),nil
}

//InstancesCount returns the total number of instances in the storage (all sessions combined)
func (s *fileStorage)InstancesCount()(int,error){
	s.rw.Lock()
	defer s.rw.Unlock()

	var total int
	for k := range s.db{
		total += len(s.db[k].Instances)
	}
	return total, nil
}

//InstanceDelete deletes an instance from the storage
func (s *fileStorage)InstanceDelete(sID, iID string)(error){
	s.rw.Lock()


	sn, ok := s.db[sID]
	if !ok{
		return fmt.Errorf("Session %s not found", sID)
	}
	delete(sn.Instances, iID)
	s.rw.Unlock()
	return s.SessionSave(sn)
}

//InstanceList returns all instances in a session in storage
func (s *fileStorage)InstanceList(sID string)([]*types.Instance, error){
	s.rw.Lock()
	defer s.rw.Unlock()

	sn, ok := s.db[sID]
	if !ok{
		return nil, fmt.Errorf("Session %s not found", sID)
	}

	var ret []*types.Instance

	for _,v := range sn.Instances{
		ret = append(ret, v)
	}

	return ret, nil
}
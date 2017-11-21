package storage

import (
	"github.com/sofianinho/vnf-api-golang/vnf/types"
)

//API serves the storage service for the different parts that it handles: session, config and instance. Typical operations are save, search, get, count, and delete
type API interface{
	SessionGet(id string)(*types.Session,error)
	SessionSave(*types.Session)(error)
	SessionCount()(int,error)
	SessionDelete(id string)(error)
	SessionList()([]*types.Session, error)

	ConfigGet(sID, cID string)(*types.Config,error)
	ConfigFindByAlias(sID, alias string)(*types.Config, error)
	ConfigFindByTags(sID string, tags []string)([]*types.Config, error)
	//ConfigAddInstance(sID, cID string, i *types.Instance)(error)
	//ConfigDelInstance(sID, cID string, i *types.Instance)(error)
	ConfigSave(sID string, c *types.Config)(error)
	ConfigDelete(sID, cID string)(error)
	ConfigCount(sID string)(int,error)
	ConfigsCount()(int, error)
	ConfigList(sID string)([]*types.Config, error)

	InstanceGet(sID, iID string)(*types.Instance,error)
	InstanceFindByAlias(sID, alias string)(*types.Instance, error)
	InstanceFindByTags(sID string, tags []string)([]*types.Instance, error)
	InstanceSave(sID string, i *types.Instance)(error)
	InstanceDelete(sID, iID string)(error)
	InstanceCount(sID string)(int,error)
	InstancesCount()(int, error)
	InstanceList(sID string)([]*types.Instance, error)
}
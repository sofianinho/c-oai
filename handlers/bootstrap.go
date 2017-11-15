package handlers

import (
	"fmt"

	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/scheduler"
	"github.com/sofianinho/vnf-api-golang/storage"
	"github.com/sofianinho/vnf-api-golang/vnf"

)

var core vnf.API

//Bootstrap initializes the new APIs
func Bootstrap() (error){
	serverHost = config.Params.GetString("server.host")
	//set the scheduler
	var s scheduler.API
	if config.Params.GetString("scheduler.type") == "docker"{
		return fmt.Errorf("Docker scheduler not supported yet...")
	} 
	//basically the only choice of scheduler for now
	var e error
	s, e = scheduler.New()
	if e != nil{
		return fmt.Errorf("error scheduler: %s", e)
	}
	//set the storage
	var t storage.API
	if config.Params.GetString("storage.type") == "postgres"{
		return fmt.Errorf("Postgres not supported yet...")
	}
	//basically the only storage option for now
	t, e = storage.New(config.Params.GetString("storage.file.path"))
	if e!=nil{
		return fmt.Errorf("error storage: %s", e)	
	}

	core = vnf.New(s, t)
	return nil
}
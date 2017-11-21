package handlers

import (
	"os"
	"fmt"

	"github.com/sofianinho/vnf-api-golang/utils"
	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/scheduler"
	"github.com/sofianinho/vnf-api-golang/storage"
	"github.com/sofianinho/vnf-api-golang/templates"
	"github.com/sofianinho/vnf-api-golang/vnf"

)

var core vnf.API
var tpl templates.API


//Bootstrap initializes the new APIs
func Bootstrap() (error){
	serverHost = config.Params.GetString("server.host")
	//set the scheduler
	var s scheduler.API
	if config.Params.GetString("scheduler.type") == "docker"{
		return fmt.Errorf("docker scheduler not supported yet")
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
		return fmt.Errorf("postgres not supported yet")
	}
	//basically the only storage option for now
	t, e = storage.New(config.Params.GetString("storage.file.path"))
	if e!=nil{
		return fmt.Errorf("error storage: %s", e)	
	}

	core = vnf.New(s, t)
	//setup the templating interface	
	tpl, e = templates.New(config.Params.GetString("templates.path"))
	if e != nil{
		return fmt.Errorf("could not init templating system: %s", e)
	}

	//setup the path for the binaries of lte-m and vanilla oai
	wd, e := os.Getwd()
	if e != nil{
		return fmt.Errorf("Unable to get current dir: %s", e)
	}
	if e:=utils.PrefixToPath(wd+"/bin/oai"); e != nil{
		return fmt.Errorf("Unable to update path: %s", e)
	}
	if e:=utils.PrefixToPath(wd+"/bin/lte-m"); e != nil{
		return fmt.Errorf("Unable to update path: %s", e)
	}
	config.Log.Debugf("Current PATH var is: %s", os.Getenv("PATH"))

	return nil
}
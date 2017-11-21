package handlers

import(
	"os"
	"net/http"
	"path/filepath"

	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/vnf"


	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//PostTask define the swagger entry 
var PostTask = endpoint.New("post", "/session/{session_id}/instance", "Create a new VNF instance in your session",
	endpoint.Path("session_id", "string", "session id for this instance", true),
	endpoint.Query("vnf_name", "string", "the VNF name for this session", false),
	endpoint.Handler(NewTask),
	endpoint.Body(instanceJson{}, "Instance to be created in the session", true),
	endpoint.Response(http.StatusCreated, types.Instance{}, "Successfully created an instance"),
	endpoint.Description("Instance creation in your session"),
	endpoint.Tags("Instance"),
)

//NewTask is an HTTP POST handler to create a new VNF instance for the session
func NewTask(c *gin.Context){
	t := c.Query("vnf_name")
	if t != "" && t != "oai" && t != "lte-m"{
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, apiReply{Status: unkVNF})
		return
	}
	//if t is empty revert to default
	if t == ""{
		t = vnf.DefaultVNFName
	}

	sessionID := c.Param("session_id")
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}
	var it instanceJson
	if c.BindJSON(&it) != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, apiReply{SessionID: sessionID, Hostname: serverHost, Status: wrongConfig})
		return
	}
	conf, err := core.ConfigGet(sessionID,it.ConfigID)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: serverHost, Status: confNotFound})
		return
	}
	ins, err := core.InstanceNew(sessionID, conf, it.Alias, it.Tags, t)
	if err != nil{
		config.Log.Errorf("Failed to create a new instance: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	//Running the task
	//1. Create a folder in runtime with the ID of the task
	taskPath := filepath.Join(config.Params.GetString("runtime.path"), ins.ID)
	if e:=os.Mkdir(taskPath, os.ModeDir|os.FileMode(0755)); e!=nil{
		//if this fails, do not forget to cleanup
		config.Log.Errorf("Cannot create task runtime dir %s: %s", taskPath, e)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		if err := core.InstanceDelete(sessionID, ins.ID); err!=nil{
			config.Log.Errorf("Failed to delete task %s of session %s from storage: %s", ins.ID, sessionID, err)
		}
		return 
	}
	//2. Compile the configuration into this folder
	if e:= tpl.CompileConfig(conf, taskPath); e!=nil{
		//if this fails, do not forget to cleanup
		config.Log.Errorf("Cannot compile confid %s into path %s: %s", conf.ID, taskPath, e)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: wrongConfig})
		if err := core.InstanceDelete(sessionID, ins.ID); err!=nil{
			config.Log.Errorf("Failed to delete task %s of session %s from storage: %s", ins.ID, sessionID, err)
		}
		return 
	}
	
	if e := core.InstanceRun(ins); e!=nil{
		config.Log.Errorf("Cannot run task %s: %s", ins.ID, e)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		if err := core.InstanceDelete(sessionID, ins.ID); err!=nil{
			config.Log.Errorf("Failed to delete task %s of session %s from storage: %s", ins.ID, sessionID, err)
		}
		return 
	}

	c.JSON(http.StatusCreated, ins)
	return 
}
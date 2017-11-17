package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)

//GetTask define the swagger entry 
var GetTask = endpoint.New("get", "/session/{session_id}/instance/{instance_id}", "Returns an instance from your VNF manager session",
	endpoint.Path("session_id", "string", "session id to search", true),
	endpoint.Path("instance_id", "string", "instance id to search", true),
	endpoint.Handler(ListTask),
	endpoint.Response(http.StatusOK, types.Instance{}, "Successfully found an instance"),
	endpoint.Description("Instance search in your VNF manager session"),
	endpoint.Tags("Instance"),
)

//GetTasks define the swagger entry 
var GetTasks = endpoint.New("get", "/session/{session_id}/instances", "Returns all instances in your VNF manager session",
	endpoint.Path("session_id", "string", "session id to search", true),
	endpoint.Handler(ListTasks),
	endpoint.Response(http.StatusOK, []types.Instance{}, "Successfully returned instances"),
	endpoint.Description("Instances listing in your VNF manager session"),
	endpoint.Tags("Instance"),
)

//ListTask is an HTTP POST handler to return an instance
func ListTask(c *gin.Context){
	sessionID := c.Param("session_id")
	instanceID := c.Param("instance_id")
	
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}

	i, err := core.InstanceGet(sessionID, instanceID)
	if err != nil{
		config.Log.Debugf("Instance search error: %s", err)
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: insNotFound})
		return
	}

	c.JSON(http.StatusOK, i)
	return
}

//ListTasks is an HTTP POST handler to return instances in a session
func ListTasks(c *gin.Context){
	sessionID := c.Param("session_id")

	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}

	ret, err := core.InstanceList(sessionID)
	if err != nil{
		config.Log.Errorf("Failed to return a list of instances in a session: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{Hostname: serverHost, Status: unkError})
		return
	}

	c.JSON(http.StatusOK, ret)
	return
}
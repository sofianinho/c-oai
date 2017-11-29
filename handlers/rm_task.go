package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//DeleteTask define the swagger entry 
var DeleteTask = endpoint.New("delete", "/session/{session_id}/instance/{instance_id}", "Delete a VNF instance in your session",
	endpoint.Path("session_id", "string", "session id for this config", true),
	endpoint.Path("instance_id", "string", "instance id to delete", true),
	endpoint.Handler(RmInstance),
	endpoint.Response(http.StatusOK, apiReply{}, "Successfully deleted an instance"),
	endpoint.Description("Instance delete in your session"),
	endpoint.Tags("Instance"),
)

//RmInstance is an HTTP DELETE handler to delete an instance for the VNF
func RmInstance(c *gin.Context){
	sessionID := c.Param("session_id")
	instanceID := c.Param("instance_id")
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}
	
	if _, err := core.InstanceGet(sessionID, instanceID); err != nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: serverHost, Status: insNotFound})
		return
	}
	
	err := core.InstanceDelete(sessionID, instanceID)
	if err != nil{
		config.Log.Errorf("Failed to delete an instance: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	
	c.JSON(http.StatusOK, apiReply{SessionID: sessionID, Hostname: serverHost, Status: insDeleted})
	return 
}
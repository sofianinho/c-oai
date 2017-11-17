package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//DeleteConfig define the swagger entry 
var DeleteConfig = endpoint.New("delete", "/session/{session_id}/config/{config_id}", "Delete a VNF configuration in your session",
	endpoint.Path("session_id", "string", "session id for this config", true),
	endpoint.Path("config_id", "string", "config id to delete", true),
	endpoint.Handler(RmConfig),
	endpoint.Response(http.StatusCreated, apiReply{}, "Successfully deleted a configuration"),
	endpoint.Description("Configuration delete in your session"),
	endpoint.Tags("Configuration"),
)

//RmConfig is an HTTP DELETE handler to delete a configuration for the VNF
func RmConfig(c *gin.Context){
	sessionID := c.Param("session_id")
	configID := c.Param("config_id")
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}
	
	if _, err := core.ConfigGet(sessionID, configID); err != nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: serverHost, Status: confNotFound})
		return
	}
	
	err := core.ConfigDelete(sessionID, configID)
	if err != nil{
		config.Log.Errorf("Failed to delete a config: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	
	c.JSON(http.StatusOK, apiReply{SessionID: sessionID, Hostname: serverHost, Status: confDeleted})
	return 
}
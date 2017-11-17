package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)

//GetConfig define the swagger entry 
var GetConfig = endpoint.New("get", "/session/{session_id}/config/{config_id}", "Returns a configuration from your VNF manager session",
	endpoint.Path("session_id", "string", "session id to search", true),
	endpoint.Path("config_id", "string", "config id to search", true),
	endpoint.Handler(ListConfig),
	endpoint.Response(http.StatusOK, types.Config{}, "Successfully found a config"),
	endpoint.Description("Config search in your VNF manager session"),
	endpoint.Tags("Configuration"),
)

//GetConfigs define the swagger entry 
var GetConfigs = endpoint.New("get", "/session/{session_id}/configs", "Returns all configs in your VNF manager session",
	endpoint.Path("session_id", "string", "session id to search", true),
	endpoint.Handler(ListConfigs),
	endpoint.Response(http.StatusOK, []types.Config{}, "Successfully returned  configs"),
	endpoint.Description("Configs listing in your VNF manager session"),
	endpoint.Tags("Configuration"),
)

//ListConfig is an HTTP POST handler to return a config
func ListConfig(c *gin.Context){
	sessionID := c.Param("session_id")
	configID := c.Param("config_id")
	
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}

	cf, err := core.ConfigGet(sessionID, configID)
	if err != nil{
		config.Log.Debugf("Config search error: %s", err)
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: confNotFound})
		return
	}

	c.JSON(http.StatusOK, cf)
	return
}

//ListConfigs is an HTTP POST handler to return configs in a session
func ListConfigs(c *gin.Context){
	sessionID := c.Param("session_id")

	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}

	ret, err := core.ConfigList(sessionID)
	if err != nil{
		config.Log.Errorf("Failed to return a list of configs in a session: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{Hostname: serverHost, Status: unkError})
		return
	}

	c.JSON(http.StatusOK, ret)
	return
}
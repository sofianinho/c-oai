package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/vnf/types"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//define the swagger entry 
var PostConfig = endpoint.New("post", "/session/{session_id}/config", "Create a new VNF configuration in your session",
	endpoint.Path("session_id", "string", "session id for this config", true),
	endpoint.Handler(NewConfig),
	endpoint.Body(confJson{}, "Configuration needs to be created in the session", true),
	endpoint.Response(http.StatusCreated, types.Config{}, "Successfully created a configuration"),
	endpoint.Description("Configuration creation in your session"),
	endpoint.Tags("Configuration"),
)

//NewConfig is an HTTP POST handler to create a new configuration for the VNF
func NewConfig(c *gin.Context){
	sessionID := c.Param("session_id")
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}
	var cf  confJson
	if c.BindJSON(&cf) != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, apiReply{SessionID: sessionID, Hostname: serverHost, Status: wrongConfig})
		return
	}
	conf, err := core.ConfigNew(s.ID, cf.Params, cf.Alias, cf.Tags)
	if err != nil{
		config.Log.Errorf("Failed to create a new config: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	c.JSON(http.StatusCreated, conf)
	return 
}
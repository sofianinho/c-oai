package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/utils"

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
	//check for asked version of template
	if tpl.VersionExists(cf.Version) != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, apiReply{SessionID: sessionID, Hostname: serverHost, Status: wrongConfigVersion})
		return
	}
	//Get the interface and IP to get the MME IP and fill in the cf.IF part
	var ifaceConf types.EnbIF
	var err error
	ifaceConf.S1MmeAddr, _, err = utils.GetRouteAndInterface(cf.Params.Enb.Mme.Ipv4)
	if err != nil{
		config.Log.Errorf("MME IPv4 error on system: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: wrongConfig})
		return
	}
	ifaceConf.S1MmeIF, err = utils.GetNameFromIfIP(ifaceConf.S1MmeAddr)
	if err != nil{
		config.Log.Errorf("Error retrieving interface name from IPv4 on system: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: wrongConfig})
		return
	}
	ifaceConf.S1UAddr = ifaceConf.S1MmeAddr
	ifaceConf.S1UIF = ifaceConf.S1MmeIF
	ifaceConf.S1UPort = 2152
	cf.Params.Enb.IF = ifaceConf
	//store the new config
	conf, err := core.ConfigNew(s.ID, cf.Version, cf.Params, cf.Alias, cf.Tags)
	if err != nil{
		config.Log.Errorf("Failed to create a new config: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	c.JSON(http.StatusCreated, conf)
	return 
}
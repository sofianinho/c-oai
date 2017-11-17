package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/vnf/types"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//PostTask define the swagger entry 
var PostTask = endpoint.New("post", "/session/{session_id}/instance", "Create a new VNF instance in your session",
	endpoint.Path("session_id", "string", "session id for this instance", true),
	endpoint.Handler(NewTask),
	endpoint.Body(instanceJson{}, "Instance to be created in the session", true),
	endpoint.Response(http.StatusCreated, types.Instance{}, "Successfully created an instance"),
	endpoint.Description("Instance creation in your session"),
	endpoint.Tags("Instance"),
)

//NewTask is an HTTP POST handler to create a new VNF instance for the session
func NewTask(c *gin.Context){
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
	ins, err := core.InstanceNew(sessionID, conf, it.Alias, it.Tags)
	if err != nil{
		config.Log.Errorf("Failed to create a new instance: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	c.JSON(http.StatusCreated, ins)
	return 
}
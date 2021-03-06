package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//PostSession define the swagger entry 
var PostSession = endpoint.New("post", "/session/", "Create a new session on your VNF manager",
//	endpoint.Query("scheduler", "string", "the VNF scheduler for this session (system, docker)", false),
	endpoint.Handler(NewSession),
	endpoint.Response(http.StatusCreated, types.Session{}, "Successfully created a new session"),
	endpoint.Description("Session creation in your VNF manager"),
	endpoint.Tags("Session"),
)

//NewSession is an HTTP POST handler to create a new configuration for the VNF
func NewSession(c *gin.Context){

/* 	sch := c.Query("scheduler")
	if sch != "" && sch != "system"{
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, apiReply{Status: wrongConfig})
		return
	} */
	s, err := core.SessionNew()
	if err != nil{
		config.Log.Errorf("Failed to create new session: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{Status: unkError})
		return
	}
	config.Log.Debugf("Created a new session: %v", s)
	c.JSON(http.StatusCreated, apiReply{SessionID: s.ID, Hostname: config.Params.GetString("server.host"), Status: sessionCreated})
	return
}
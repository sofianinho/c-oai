package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//define the swagger entry 
var PostSession = endpoint.New("post", "/session/", "Create a new session on your VNF manager",
	endpoint.Query("vnf_name", "string", "the VNF name for this session", false),
	endpoint.Handler(NewSession),
	endpoint.Response(http.StatusCreated, types.Session{}, "Successfully created a new session"),
	endpoint.Description("Session creation in your VNF manager"),
	endpoint.Tags("Session"),
)

//NewSession is an HTTP POST handler to create a new configuration for the VNF
func NewSession(c *gin.Context){
	vnfName := c.Query("vnf_name")
	if vnfName != "" && vnfName != "oai"{
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, apiReply{Status: unkVNF})
		return
	}
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
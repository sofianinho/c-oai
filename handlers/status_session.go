package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//GetStatusSession define the swagger entry 
var GetStatusSession = endpoint.New("get", "/status", "Returns a status for the scheduler",
	endpoint.Handler(StatusSession),
	endpoint.Response(http.StatusOK, types.Status{}, "Successfully returned a system status"),
	endpoint.Description("Scheduler status in your VNF manager"),
	endpoint.Tags("Status"),
)

//StatusSession is an HTTP GET handler to return a session status for the VNF manager
func StatusSession(c *gin.Context){

	st, err := core.SessionStatus()
	if err != nil{
		config.Log.Errorf("Failed to retrieve system status: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{Status: unkError})
		return
	}
	
	c.JSON(http.StatusOK, st)
	return
}
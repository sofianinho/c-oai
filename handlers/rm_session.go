package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)


//DeleteSession define the swagger entry 
var DeleteSession = endpoint.New("delete", "/session/{session_id}", "Delete a session",
	endpoint.Path("session_id", "string", "session id to delete", true),
	endpoint.Handler(RmSession),
	endpoint.Response(http.StatusOK, apiReply{}, "Successfully deleted a session"),
	endpoint.Description("Delete a session with configuration and instances"),
	endpoint.Tags("Session"),
)

//RmSession is an HTTP DELETE handler to delete a session
func RmSession(c *gin.Context){
	sessionID := c.Param("session_id")
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}
	
	err := core.SessionDelete(sessionID)
	if err != nil{
		config.Log.Errorf("Failed to delete a session: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{SessionID: sessionID, Hostname: serverHost, Status: unkError})
		return
	}
	
	c.JSON(http.StatusOK, apiReply{SessionID: sessionID, Hostname: serverHost, Status: sessionClosed})
	return 
}
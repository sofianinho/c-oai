package handlers

import(
	"net/http"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/config"

	"github.com/savaki/swag/endpoint"
	"github.com/gin-gonic/gin"
)

//GetSession define the swagger entry 
var GetSession = endpoint.New("get", "/session/{session_id}", "Returns a session on your VNF manager",
	endpoint.Path("session_id", "string", "session id to search", true),
	endpoint.Handler(ListSession),
	endpoint.Response(http.StatusOK, types.Session{}, "Successfully found a session"),
	endpoint.Description("Session search in your VNF manager"),
	endpoint.Tags("Session"),
)

//GetSessions define the swagger entry 
var GetSessions = endpoint.New("get", "/sessions/", "Returns all sessions on your VNF manager",
	endpoint.Handler(ListSessions),
	endpoint.Response(http.StatusOK, []types.Session{}, "Successfully returned  sessions"),
	endpoint.Description("Sessions listing in your VNF manager"),
	endpoint.Tags("Session"),
)

//ListSession is an HTTP POST handler to return a session
func ListSession(c *gin.Context){
	sessionID := c.Param("session_id")
	
	s := core.SessionGet(sessionID)
	if s == nil{
		c.AbortWithStatusJSON(http.StatusNotFound, apiReply{SessionID: sessionID, Hostname: config.Params.GetString("server.host"), Status: sessionNotFound})
		return
	}

	c.JSON(http.StatusOK, s)
	return
}

//ListSessions is an HTTP POST handler to return sessions
func ListSessions(c *gin.Context){
	ret, err := core.SessionList()
	if err != nil{
		config.Log.Errorf("Failed to return a list of sessions: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, apiReply{Hostname: serverHost, Status: unkError})
		return
	}

	c.JSON(http.StatusOK, ret)
	return
}
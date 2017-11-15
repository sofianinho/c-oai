package handlers

import(
	"fmt"
	"github.com/sofianinho/vnf-api-golang/config"
)
const (
	unkError = "Server unexpected failure. Try again later."
	unkVNF = "Unknown or unsupported vnf"
	sessionCreated = "Session created"
	sessionNotFound = "Session not found"
	sessionClosed = "Session closed and cleaned up"
	insNotFound = "Instance not found"
	insCreated = "Instance created"
	insDeleted = "Instance deleted and cleaned up"
	confNotFound = "Config not found"
	confCreated = "Config created"
	confDeleted = "Config deleted and cleaned up"
)

var serverHost string

type apiReply struct{
	SessionID string `json:"session_id,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	Status	  string `json:"status"`
}

var urlPath = fmt.Sprintf("%s%s", config.ApiSubpath, config.ApiCurrentVersion)

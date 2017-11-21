package types

import (
	"time"
)

//Instance is for starting an instance of the VNF with its configuration
type Instance struct {
	ID			string		`json:"id,omitempty"`
	Alias		string		`json:"alias,omitempty"`
	Session		string		`json:"session_id"`
	CreatedAt	time.Time   `json:"created_at"`
	Tags		[]string	`json:"tags,omitempty"`
	Artefact	string
	Config		*Config		`json:"configuration"`
}

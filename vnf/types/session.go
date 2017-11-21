package types

import (
	"time"
	"sync"
)

//Session is for keeping track of the session: configurations, instances and other
type Session struct {
	ID			string					`json:"id,omitempty"`
	Instances	map[string]*Instance	`json:"instances"`
	Configs		map[string]*Config		`json:"configurations"`
	CreatedAt	time.Time   			`json:"created_at,omitempty"`
	rw			sync.Mutex
}

//Lock blocks the access to a session
func (s *Session) Lock(){
	s.rw.Lock()
}

//Unlock unblocks the access to a session
func (s *Session) Unlock(){
	s.rw.Unlock()
}
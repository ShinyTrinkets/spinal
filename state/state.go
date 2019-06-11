package state

import (
	"time"

	ovr "github.com/ShinyTrinkets/overseer.go"
)

type stateTree = map[string]Level1

// Level1 is ...
type Level1 struct {
	Props    Header1           `json:"props"`
	Children map[string]Level2 `json:"children"`
}

// Header1 represents Level1 properties
type Header1 struct {
	Enabled bool      `json:"enabled"`
	ID      string    `json:"id"`
	Db      bool      `json:"db,omitempty"`
	Log     bool      `json:"log,omitempty"`
	Path    string    `json:"path"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// Level2 is ...
type Level2 struct {
	Props Header2 `json:"props"`
}

// Header2 represents Level2 properties
type Header2 = ovr.JSONProcess

var state = stateTree{}

// GetState returns a copy of the state
func GetState() stateTree {
	return state
}

// HasLevel1 checks for a lvl1 name
func HasLevel1(name string) (exists bool) {
	_, exists = state[name]
	return
}

// GetLevel1 returns a lvl1 state
func GetLevel1(name string) Level1 {
	return state[name]
}

// SetLevel1 updates the StateTree
func SetLevel1(name string, props *Header1) {
	children := map[string]Level2{}
	state[name] = Level1{*props, children}
}

// SetLevel2 updates the StateTree
func SetLevel2(name1 string, name2 string, props *Header2) {
	state[name1].Children[name2] = Level2{*props}
}

package state

import (
	"sync"
	"time"

	ovr "github.com/ShinyTrinkets/overseer"
)

const separator = "∷"

// There is only 1 state tree and cannot be changed
var state sync.Map

// Header1 represents Level1 properties
type Header1 struct {
	Enabled bool      `json:"enabled"`
	ID      string    `json:"id"`
	Db      bool      `json:"db,omitempty"`
	Log     bool      `json:"log,omitempty"`
	Cwd     string    `json:"cwd,omitempty"`
	Path    string    `json:"path"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// Header2 represents Level2 properties
type Header2 = ovr.ProcessJSON

// GetState returns a full copy of the state
func GetState() sync.Map {
	return state
}

// HasLevel1 checks for a lvl1 name
func HasLevel1(name string) (exists bool) {
	_, exists = state.Load(name)
	return
}

// GetLevel1 returns a lvl1 state
// Level1 represents one recipe/file
func GetLevel1(name string) Header1 {
	l, _ := state.Load(name)
	return l.(Header1)
}

// SetLevel1 updates the StateTree
func SetLevel1(name string, props *Header1) {
	state.Store(name, *props)
}

// HasLevel2 checks for a lvl2 name
func HasLevel2(name1 string, name2 string) (exists bool) {
	_, exists = state.Load(name1 + separator + name2)
	return
}

// GetLevel2 returns a lvl2 state
// Level2 represents one child from a recipe/file
func GetLevel2(name1 string, name2 string) Header2 {
	l, _ := state.Load(name1 + separator + name2)
	return l.(Header2)
}

// SetLevel2 updates the StateTree
func SetLevel2(name1 string, name2 string, props *Header2) {
	state.Store(name1+separator+name2, *props)
}

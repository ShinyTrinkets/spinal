package state

import (
	"sync"
	"time"

	ovr "github.com/ShinyTrinkets/overseer"
)

// type stateTree = map[string]Level1

// Level1 represents one recipe/file
// the children are the code files included in it
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

// Level2 represents one child from a recipe/file
type Level2 struct {
	Props Header2 `json:"props"`
}

// Header2 represents Level2 properties
type Header2 = ovr.ProcessJSON

// There is only 1 state tree and cannot be changed
var state sync.Map

// GetState returns a copy of the state
func GetState() sync.Map {
	return state
}

// HasLevel1 checks for a lvl1 name
func HasLevel1(name string) (exists bool) {
	_, exists = state.Load(name)
	return
}

// GetLevel1 returns a lvl1 state
func GetLevel1(name string) Level1 {
	l, _ := state.Load(name)
	return l.(Level1)
}

// SetLevel1 updates the StateTree
func SetLevel1(name string, props *Header1) {
	children := map[string]Level2{}
	state.Store(name, Level1{*props, children})
}

// HasLevel2 checks for a lvl2 name
func HasLevel2(name1 string, name2 string) (exists bool) {
	l := GetLevel1(name1)
	_, exists = l.Children[name2]
	return
}

// GetLevel2 returns a lvl2 state
func GetLevel2(name1 string, name2 string) Level2 {
	l1 := GetLevel1(name1)
	return l1.Children[name2]
}

// SetLevel2 updates the StateTree
func SetLevel2(name1 string, name2 string, props *Header2) {
	l := GetLevel1(name1)
	l.Children[name2] = Level2{*props}
}

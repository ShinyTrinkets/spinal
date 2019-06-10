package state

import (
	"fmt"
	"time"

	ovr "github.com/ShinyTrinkets/overseer.go"
)

// State is
type State = map[string]Level1

// Level1 is ...
type Level1 struct {
	Props    Header1           `json:"props"`
	Children map[string]Level2 `json:"children"`
}

// Header1 is ...
type Header1 struct {
	Enabled bool      `json:"spinal"`
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

// Header2 is ...
type Header2 = ovr.JSONProcess

var state = State{}

// GetState returns a copy of the state
func GetState() State {
	return state
}

// SetLevel1 is ...
func SetLevel1(name string, props Header1) {
	children := map[string]Level2{}
	state[name] = Level1{props, children}
	fmt.Println(state[name])
	fmt.Println("====")
}

// SetLevel2 is ...
func SetLevel2(name1 string, name2 string, props Header2) {
	state[name1].Children[name2] = Level2{props}
	fmt.Println(state[name1].Children[name2])
	fmt.Println("====")
}

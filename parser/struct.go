package parser

import (
	"time"
)

const (
	validSourceExt = ".md"
	maxFolderDepth = 3
	blankRunes     = "\t\n\r "
)

type StringToString map[string]string

type FrontMatter struct {
	Enabled bool   `yaml:"spinal"`
	Id      string `yaml:"id"`
	Db      bool   `yaml:"db,omitempty"`
	Log     bool   `yaml:"log,omitempty"`
}

type CodeFile struct {
	FrontMatter
	Path   string
	Ctime  time.Time
	Mtime  time.Time
	Blocks map[string]string
}

// Minimal validation check for ID and Path
func (self *CodeFile) IsValid() bool {
	if len(self.Path) < 2 {
		return false
	}
	l := len(self.Id)
	return (l > 0 && l < 100)
}

type CodeType struct {
	Name       string
	Executable string
	Comment    string
}

// All known code block types
var CodeBlocks = map[string]CodeType{
	"js": {"Javascript", "node", "//"},
	"py": {"Python", "python", "#"},
	// "rb": {"Ruby", "ruby", "#"},
}

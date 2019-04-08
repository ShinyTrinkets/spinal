package parser

import (
	"time"
)

const (
	validSourceExt  = ".md"
	maxSrcScanDepth = 3
	blankRunes      = "\t\n\r "
)

type StringToString map[string]string

type MetaData interface{}

type FrontMatter struct {
	Enabled    bool     `yaml:"spinal" json:"spinal"`
	Id         string   `yaml:"id"  json:"id"`
	Db         bool     `yaml:"db,omitempty"  json:"db,omitempty"`
	Log        bool     `yaml:"log,omitempty" json:"log,omitempty"`
	DelayStart uint     `yaml:"delayStart,omitempty" json:"delayStart,omitempty"`
	RetryTimes uint     `yaml:"retryTimes,omitempty" json:"retryTimes,omitempty"`
	Meta       MetaData `yaml:"meta" json:"meta"`
}

type CodeFile struct {
	FrontMatter
	Path   string
	Ctime  time.Time
	Mtime  time.Time
	Blocks map[string]string
}

// IsValid makes a validation check for ID and Path
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
	"py": {"Python", "python3", "#"},
	// "rb": {"Ruby", "ruby", "#"},
}

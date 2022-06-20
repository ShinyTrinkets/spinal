package state

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stateLength(s *sync.Map) (i int) {
	s.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return
}

func TestStateLvls(t *testing.T) {
	assert := assert.New(t)

	s := GetState()
	assert.Equal(0, stateLength(&s))

	SetLevel1("x.md",
		&Header1{
			Enabled: true,
			ID:      "x",
			Db:      true,
			Log:     true,
			Path:    "x/y/z",
		})

	s = GetState()
	assert.Equal(1, stateLength(&s))

	assert.True(HasLevel1("x.md"))
	assert.True(GetLevel1("x.md").Enabled)
	assert.Equal("x", GetLevel1("x.md").ID)
	assert.Equal("x/y/z", GetLevel1("x.md").Path)

	SetLevel2("x.md", "x.js",
		&Header2{
			ID:  "x",
			Cmd: "x.js",
			Dir: ".",
		})

	s = GetState()
	assert.Equal(2, stateLength(&s))

	assert.True(HasLevel2("x.md", "x.js"))
	assert.Equal("x", GetLevel2("x.md", "x.js").ID)
	assert.Equal(".", GetLevel2("x.md", "x.js").Dir)
}

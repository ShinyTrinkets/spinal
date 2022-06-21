package kvstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeUnit = 100 * time.Millisecond

func TestTable(t *testing.T) {
	assert := assert.New(t)

	table := CacheTable{name: "table1"}
	assert.Equal(0, table.Count())
	assert.False(table.Exists("x"))

	table.Set("x", 123, timeUnit)
	assert.Equal(1, table.Count())
	assert.True(table.Exists("x"))
	v, ok := table.Get("x")
	assert.True(ok)
	assert.Equal(123, v)

	table.Set("x", "XYZ", timeUnit)
	assert.Equal(1, table.Count())
	assert.True(table.Exists("x"))
	v, ok = table.Get("x")
	assert.True(ok)
	assert.Equal("XYZ", v)

	table.StartCleaner(timeUnit)
	time.Sleep(timeUnit)

	assert.Equal(0, table.Count())
	assert.False(table.Exists("x"))
	v, ok = table.Get("x")
	assert.False(ok)
	assert.Nil(v)
}

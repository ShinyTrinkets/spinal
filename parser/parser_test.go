package parser

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type Fixtures struct {
	Text   string            `yaml:"text"`
	Result map[string]string `yaml:"result"`
}

func TestParseHeader(t *testing.T) {
	assert := assert.New(t)
	var text string
	// Without extra spaces
	text = "---\nid: 1\n---\nblah blah\n\n"
	h, b := splitHeadBody(text)
	assert.Equal(h, "---\nid: 1\n---", "Parse headers invalid header")
	assert.Equal(b, "blah blah", "Parse headers invalid body")
	// With extra spaces
	text = "---\n\n\nid:    1\n\n---\n\n\nblah blah\n\n"
	h, b = splitHeadBody(text)
	assert.Equal(h, "---\n\n\nid:    1\n\n---", "Parse headers invalid header")
	assert.Equal(b, "blah blah", "Parse headers invalid body")
}

func TestParseBlocks(t *testing.T) {
	// Testing block extraction
	text, err := ioutil.ReadFile("testdata/parse-blocks.yml")
	if err != nil {
		t.Fatalf("Cannot open YAML fixtures file")
	}
	fixtures := []Fixtures{}
	if err := yaml.Unmarshal(text, &fixtures); err != nil {
		t.Fatalf("Cannot parse YAML fixtures")
	}

	for _, fixt := range fixtures {
		blocks := ParseBlocks(fixt.Text)
		for lang, code := range fixt.Result {
			code = strings.Trim(code, blankRunes)
			bloc := strings.Trim(blocks[lang], blankRunes)
			if bloc != code {
				t.Fatalf("Resulted block = `%v` invalid ; expected = `%v`", bloc, code)
			}
		}
	}
}

func TestListCodeFiles(t *testing.T) {
	assert := assert.New(t)
	// Testing listing code files, depth 1
	files, err := listCodeFiles("testdata/deep1/", 1)
	assert.Nil(err, "Cannot list code files")
	srcFiles, err := filepath.Glob("testdata/deep1/*.md")
	assert.Nil(err, "Cannot glob code files")

	assert.Equal(len(files), len(srcFiles), "There should be %v code files != %v", len(srcFiles), len(files))

	// Testing depth 2
	files, err = listCodeFiles("testdata/deep1/", 2)
	assert.Equal(len(files), len(srcFiles)+1, "There should be %v code files != %v", len(srcFiles)+1, len(files))
	assert.Nil(err)
	// Testing depth 3
	files, err = listCodeFiles("testdata/deep1/", 3)
	assert.Equal(len(files), len(srcFiles)+2, "There should be %v code files != %v", len(srcFiles)+2, len(files))
	assert.Nil(err)
}

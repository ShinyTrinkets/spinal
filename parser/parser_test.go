package parser

import (
	"io/ioutil"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

type Fixtures struct {
	Text   string            `yaml:"text"`
	Result map[string]string `yaml:"result"`
}

func TestParseHeader(t *testing.T) {
	var text string
	// Without extra spaces
	text = "---\nid:  1\n---\nblah blah\n\n"
	h, b := splitHeadBody(text)
	if h != "---\nid:  1\n---" {
		t.Fatalf("Parse headers invalid header = %v", h)
	}
	if b != "blah blah" {
		t.Fatalf("Parse headers invalid body = %v", b)
	}
}

func TestParseBlocks(t *testing.T) {
	text, err := ioutil.ReadFile("fixtures/parse-blocks.yml")
	if err != nil {
		t.Fatalf("Cannot open YAML fixtures file")
	}
	fixtures := []Fixtures{}
	if err := yaml.Unmarshal(text, &fixtures); err != nil {
		t.Fatalf("Cannot parse YAML fixtures")
	}
	fm := FrontMatter{true, "x", false, true}

	for _, fixt := range fixtures {
		blocks := ParseBlocks(fm, fixt.Text)
		if len(blocks) < 1 {
			t.Fatalf("Len of blocks = %v ; want more than 0", len(blocks))
		}
		for lang, code := range fixt.Result {
			code = strings.Trim(code, blankRunes)
			bloc := strings.Trim(blocks[lang], blankRunes)
			if bloc != code {
				t.Fatalf("Resulted block = `%v` invalid ; expected = `%v`", bloc, code)
			}
		}
	}
}

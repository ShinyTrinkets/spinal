package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/immortal/xtime"
	"gopkg.in/yaml.v2"
)

// Lists all candidate text files from a folder.
// Candidate files should contain fenced code blocks.
// This list can be used to parse the files,
// or generate code files from the code blocks.
func ListCodeFiles(folder string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Code files must have a valid file extension
		if isFile(path) && filepath.Ext(f.Name()) == VALID_EXT {
			fileList = append(fileList, path)
		}
		return nil
	})
	// Possible errors in case the folder cannot be accessed
	return fileList, err
}

// Parse a file and return a structure.
// The structure contains ID, Path, Creation time and
// a list of blocks of code.
// If the file can't be accessed or parsed,
// the returned structure will be empty and invalid.
func ParseFile(fname string) CodeFile {
	parseFile := CodeFile{}

	ctime, mtime, err := fileStats(fname)
	if err != nil {
		// os.Stat error => ignore file
		return parseFile
	}
	text, err := ioutil.ReadFile(fname)
	if err != nil {
		// Read file error => ignore file
		return parseFile
	}

	h, b := SplitHeadBody(string(text))

	fm := FrontMatter{}
	if err := yaml.Unmarshal([]byte(h), &fm); err != nil {
		// YAML parse error => ignore file
		return parseFile
	}

	return CodeFile{fm, fname, ctime, mtime, ParseBlocks(fm, b)}
}

// Split text file into front-header and the rest of the text
func SplitHeadBody(text string) (string, string) {
	re := regexp.MustCompile("(?sU)^---[\n\r]+.+[\n\r]+---[\n\r]")
	head := strings.TrimRight(re.FindString(text), BLANKS)
	body := strings.Trim(text[len(head):], BLANKS)
	return head, body
}

// Helper that returns true if the path is a regular file
func isFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); !m.IsDir() && m.IsRegular() && m&400 != 0 {
		return true
	}
	return false
}

// Helper that returns File stats (creation and modif times)
func fileStats(fname string) (time.Time, time.Time, error) {
	var c time.Time
	var m time.Time

	fi, err := os.Stat(fname)
	if err != nil {
		// File stats error
		return c, m, err
	}

	c = xtime.Get(fi).Ctime()
	m = xtime.Get(fi).Mtime()
	return c, m, nil
}

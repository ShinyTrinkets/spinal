package parser

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// For readability, higher level functions go first

// ConvertFolder finds all candidate text files from a folder,
// and generates code files.
// The original text files are not changed.
func ConvertFolder(dir string) (map[string]StringToString, error) {
	result := map[string]StringToString{}

	files, err := ParseFolder(dir, true)
	if err != nil {
		return result, err
	}

	for _, p := range files {
		outFiles, err := ConvertFile(p, false)
		if err != nil {
			// What should this do if the file cannot be converted ?
			continue // => silently ignore ?
		}
		result[p.Path] = outFiles
	}
	return result, nil
}

// ParseFolder finds all candidate text files from a folder,
// and extracts useful info about them.
func ParseFolder(dir string, checkInvalid bool) ([]CodeFile, error) {
	files := []CodeFile{}

	// The path must have a valid name
	if len(dir) < 2 {
		return files, errors.New("folder name too short: " + dir)
	}
	// and must exist locally
	stat, err := os.Stat(dir)
	if err != nil {
		return files, err
	}
	// and must be a folder
	if !stat.IsDir() {
		return files, errors.New("no such folder: " + dir)
	}

	filesStr, err := listCodeFiles(dir, 0)
	if err != nil {
		return files, err
	}
	if len(filesStr) < 1 {
		return files, errors.New("no candidate files found")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return files, err
	}

	for _, fname := range filesStr {
		// Parsing needs full absolute path
		p := ParseFile(fname)
		if checkInvalid {
			if !p.IsValid() {
				continue
			}
			if len(p.Blocks) < 1 {
				continue
			}
		}
		// resolve path relative to current dir
		if strings.Index(p.Path, cwd) == 0 {
			p.Path, err = filepath.Rel(cwd, p.Path)
			if err != nil {
				continue
			}
		}
		files = append(files, p)
	}
	return files, nil
}

// listCodeFiles returns all candidate text files from a folder.
// Candidate files should contain fenced code blocks.
// This list can be used to parse the files,
// or generate code files from the code blocks.
func listCodeFiles(folder string, scanDepth int) ([]string, error) {
	fileList := []string{}
	if scanDepth < 1 {
		scanDepth = maxSrcScanDepth
	}

	// Resolve absolute path
	absDir, err := filepath.Abs(folder)
	if err != nil {
		return fileList, err
	}
	// Resolve potential symlinks
	realDir, err := filepath.EvalSymlinks(absDir)
	if err != nil {
		return fileList, err
	}
	baseLen := len(realDir) + 1

	err = filepath.Walk(realDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Code files must have a valid file extension
		if isFile(path) && filepath.Ext(f.Name()) == validSourceExt {
			// Count the slashes to estimate folder depth
			if strings.Count(path[baseLen:], "/") >= scanDepth {
				return nil
			}
			// Append the full absolute path
			fileList = append(fileList, path)
		}
		return nil
	})
	// Possible errors in case the folder cannot be accessed
	return fileList, err
}

// ConvertFile generates 1 or more code files, from one code file.
func ConvertFile(codFile CodeFile, force bool) (StringToString, error) {
	outFiles := StringToString{}

	// Some logic to decide if the structure is valid
	if !force && !codFile.IsValid() {
		return outFiles, errors.New("file header is invalid: " + codFile.Path)
	}
	// The code file must be enabled
	if !force && !codFile.Enabled {
		return outFiles, errors.New("file is marked disabled: " + codFile.Path)
	}
	// And must have blocks of code
	if len(codFile.Blocks) == 0 {
		return outFiles, errors.New("file has no blocks of code: " + codFile.Path)
	}

	fName := codFile.Path
	baseLen := len(fName) - len(filepath.Ext(fName))

	front := FrontMatter{codFile.Enabled, codFile.Id, codFile.Db, codFile.Log, codFile.Meta}

	for lang, code := range codFile.Blocks {
		outFile := fName[:baseLen] + "." + lang
		code = codeGeneratedByMsg(lang) + "\n\n" +
			codeLangHeader(front, lang) + "\n" +
			codeLangImports(front, lang) + "\n" + code
		err := ioutil.WriteFile(outFile, []byte(code), 0644)
		if err != nil {
			return outFiles, err
		}
		outFiles[lang] = outFile
	} // for each block of code
	return outFiles, nil
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

	h, b := splitHeadBody(string(text))

	fm := FrontMatter{}
	if err := yaml.Unmarshal([]byte(h), &fm); err != nil {
		// YAML parse error => ignore file
		return parseFile
	}

	// Unmarshal meta data
	var meta MetaData
	if err := yaml.Unmarshal([]byte(h), &meta); err != nil {
		// YAML parse error => ignore file
		return parseFile
	}
	if meta != nil {
		fmTags := getTagsByName(fm, "json")
		fm.Meta = normalizeMapIgnore(meta, fmTags)
	}

	return CodeFile{fm, fname, ctime, mtime, ParseBlocks(fm, b)}
}

// splitHeadBody splits a text into front-header and body-the rest of the text
func splitHeadBody(text string) (string, string) {
	re := regexp.MustCompile("(?sU)^---[\n\r]+.+[\n\r]+---[\n\r]")
	head := strings.TrimRight(re.FindString(text), blankRunes)
	body := strings.Trim(text[len(head):], blankRunes)
	return head, body
}

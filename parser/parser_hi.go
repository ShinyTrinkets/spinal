package parser

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type StringToString map[string]string

// For readability, higher level functions go first

// Find all candidate text files from a folder,
// and generate code files.
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

// Find all candidate text files from a folder,
// and extract useful info about them.
func ParseFolder(dir string, checkInvalid bool) ([]CodeFile, error) {
	files := []CodeFile{}
	cwd, err := os.Getwd()
	if err != nil {
		return files, err
	}

	// The path must have a valid name
	if len(dir) < 2 {
		return files, errors.New("folder name too short: " + dir)
	}
	fi, err := os.Stat(dir)
	if err != nil {
		return files, err
	}
	if !fi.IsDir() {
		return files, errors.New("no such folder: " + dir)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return files, err
	}
	filesStr, err := ListCodeFiles(absDir)
	if err != nil {
		return files, err
	}
	if len(filesStr) < 1 {
		return files, errors.New("no candidate files found")
	}

	for _, fname := range filesStr {
		p := ParseFile(fname)
		if checkInvalid {
			if !p.IsValid() {
				continue
			}
			if len(p.Blocks) < 1 {
				continue
			}
		}
		p.Path, err = filepath.Rel(cwd, p.Path)
		if err != nil {
			continue
		}
		files = append(files, p)
	}
	return files, nil
}

// Convert a single text file into 1 or more code files.
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
	front := FrontMatter{codFile.Enabled, codFile.Id, codFile.Db, codFile.Log}

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

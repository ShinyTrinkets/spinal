package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ovr "github.com/ShinyTrinkets/overseer.go"
	srv "github.com/ShinyTrinkets/spinal/http"
	parse "github.com/ShinyTrinkets/spinal/parser"
	"github.com/ShinyTrinkets/spinal/state"
)

type (
	strToStr = parse.StringToString
	codeFile = parse.CodeFile
)

// SpinUp receives a folder, finds all valid source-files and runs them.
// The call is blocked untill all procs finish,
// or SIGINT or SIGTERM are sent to the parent process.
// Force is enabled only for files, it can be dangerous for folders.
// For dry run, the HTTP server and the Overseer will not run.
func SpinUp(fname string, force bool, httpOpts string, noHTTP bool, dryRun bool) {
	var (
		dir    string
		pairs  map[string]strToStr
		parsed map[string]codeFile
	)

	f, err := os.Stat(fname)
	if err != nil {
		fmt.Printf("Cannot spin-up! Invalid path: %s", fname)
		return
	}
	m, f := f.Mode(), nil

	// Overwrite HTTP options in case of dry-run
	if dryRun {
		noHTTP = true
	}
	if noHTTP {
		httpOpts = ""
	}

	if !m.IsDir() && m.IsRegular() && m&400 != 0 {
		// is file?
		p := parse.ParseFile(fname)
		outFiles, err := parse.ConvertFile(p, force)
		if err != nil {
			fmt.Printf("Cannot convert file! Error: %v", err)
			return
		}
		if force {
			fmt.Println("Unsafe mode enabled!")
		}

		fmt.Printf("Converting source file '%s' ...\n", fname)
		dir = filepath.Dir(fname)
		pairs = map[string]strToStr{p.Path: outFiles}
		parsed = map[string]codeFile{p.Path: p}

	} else if m.IsDir() && m&400 != 0 {
		// is folder?
		dir = strings.TrimRight(fname, "/")
		fmt.Printf("Converting all source-files from '%s' ...\n", dir)
		pairs, parsed, err = parse.ConvertFolder(dir)
		if err != nil {
			fmt.Printf("Cannot convert folder! Error: %v", err)
			return
		}

	} else {
		fmt.Printf("Cannot run! Invalid path!")
		return
	}

	o := ovr.NewOverseer()

	baseLen := len(dir) + 1
	for inFile, convFiles := range pairs {
		codeFile := parsed[inFile]
		// Update state LVL 1
		state.SetLevel1(inFile,
			&state.Header1{
				Enabled: codeFile.Enabled,
				ID:      codeFile.ID,
				Db:      codeFile.Db,
				Log:     codeFile.Log,
				Path:    codeFile.Path,
				Ctime:   codeFile.Ctime,
				Mtime:   codeFile.Mtime,
			})

		for lang, outFile := range convFiles {
			fmt.Printf("%s ==> %s\n", inFile, outFile)
			if dryRun {
				continue
			}

			exe := parse.CodeBlocks[lang].Executable
			// Register the process with Overseer
			p := o.Add(outFile, exe, outFile[baseLen:])
			p.SetDir(dir)

			env := append(os.Environ(), "SPIN_FILE="+outFile)
			p.SetEnv(env)

			p.Lock()
			if codeFile.DelayStart > 0 {
				p.DelayStart = codeFile.DelayStart
			}
			if codeFile.RetryTimes > 0 {
				p.RetryTimes = codeFile.RetryTimes
			}
			p.Unlock()
			// Update state LVL 2
			props2 := p.ToJSON()
			state.SetLevel2(inFile, outFile, &props2)
		}
	}

	if dryRun {
		fmt.Println("\nSimulation over.")
		return
	}

	go func() {
		if noHTTP || len(httpOpts) == 0 {
			fmt.Println("HTTP server disabled")
			return
		}
		// Setup HTTP server
		http := srv.NewServer(httpOpts)
		// Activate Overseer endpoints
		srv.OverseerEndpoint(http, o)
		srv.Serve(http)
	}()

	fmt.Println("Starting procs. Press Ctrl+C to stop...")
	o.SuperviseAll()
	fmt.Println("\nShutdown.")
}

package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ovr "github.com/ShinyTrinkets/overseer.go"
	srv "github.com/ShinyTrinkets/spinal/http"
	parse "github.com/ShinyTrinkets/spinal/parser"
)

// RunOne receives a source-file, converts and runs it.
// The call is blocked untill all procs finish,
// or SIGINT or SIGTERM are sent to the parent process.
func RunOne(fname string, force bool) {

	fi, err := os.Stat(fname)
	if err != nil {
		fmt.Printf("Cannot Run-One! Error: %v", err)
		return
	}
	if m := fi.Mode(); m.IsDir() || !m.IsRegular() || m&400 == 0 {
		fmt.Printf("Cannot Run-One, invalid file! Path: %s", fname)
		return
	}

	parseFile := parse.ParseFile(fname)
	convFiles, err := parse.ConvertFile(parseFile, force)
	if err != nil {
		fmt.Printf("Cannot convert file! Error: %v", err)
		return
	}
	if force {
		fmt.Println("Unsafe mode enabled!")
	}

	o := ovr.NewOverseer()

	dir := filepath.Dir(fname)
	baseLen := len(dir) + 1

	for lang, outFile := range convFiles {
		exe := parse.CodeBlocks[lang].Executable
		env := append(os.Environ(), "SPIN_FILE="+outFile)

		p := o.Add(parseFile.Id, exe, outFile[baseLen:])
		p.SetDir(dir)
		p.SetEnv(env)
		// p.SetStateListener(func(state ovr.CmdState) {
		// 	fmt.Println("Proc State Changed:", state)
		// })

		if parseFile.DelayStart > 0 {
			p.SetDelayStart(parseFile.DelayStart)
		}
		if parseFile.RetryTimes > 0 {
			p.SetRetryTimes(parseFile.RetryTimes)
		}
	}

	fmt.Println("Starting procs. Press Ctrl+C to stop...")
	o.SuperviseAll()
	fmt.Println("\nShutdown.")
}

// RunAll receives a folder, finds all valid source-files and runs them.
// The call is blocked untill all procs finish,
// or SIGINT or SIGTERM are sent to the parent process.
func RunAll(rootDir string, httpOpts string, noHTTP bool, dryRun bool) {

	dir := strings.TrimRight(rootDir, "/")
	// This function will perform all folder checks
	// (valid name, valid folder, containing source files)
	pairs, parsed, err := parse.ConvertFolder(dir)
	if err != nil {
		fmt.Printf("Cannot Run-All! Error: %v", err)
		return
	}

	// Overwrite HTTP options in case of dry-run
	if dryRun {
		noHTTP = true
		httpOpts = ""
	}
	o := ovr.NewOverseer()

	go func() {
		if noHTTP {
			fmt.Println("HTTP server disabled")
			return
		}
		// Setup HTTP server
		http := srv.NewServer(httpOpts)
		// Enable Overseer endpoints
		srv.OverseerEndpoint(http, o)
		srv.Serve(http)
	}()

	baseLen := len(dir) + 1
	for infile, convFiles := range pairs {
		for lang, outFile := range convFiles {
			fmt.Printf("%s ==> %s\n", infile, outFile)

			exe := parse.CodeBlocks[lang].Executable
			env := append(os.Environ(), "SPIN_FILE="+outFile)
			p := o.Add(outFile, exe, outFile[baseLen:])
			p.SetDir(dir)
			p.SetEnv(env)

			parseFile := parsed[infile]
			if parseFile.DelayStart > 0 {
				p.SetDelayStart(parseFile.DelayStart)
			}
			if parseFile.RetryTimes > 0 {
				p.SetRetryTimes(parseFile.RetryTimes)
			}
		}
	}

	if dryRun {
		fmt.Println("\nSimulation ovr.")
	} else {
		fmt.Println("Starting procs. Press Ctrl+C to stop...")
		o.SuperviseAll()
		fmt.Println("\nShutdown.")
	}
}

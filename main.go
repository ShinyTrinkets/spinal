package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ShinyTrinkets/overseer.go"
	"github.com/ShinyTrinkets/spinal/http"
	"github.com/ShinyTrinkets/spinal/logger"
	"github.com/ShinyTrinkets/spinal/parser"
	log "github.com/azer/logger"
	"github.com/jawher/mow.cli"
)

const (
	Name    = "Spinal"
	Descrip = "c[○┬●]כ"
)

var (
	Version    string // injected by go build
	CommitHash string
	BuildTime  string
)

var (
	dbg bool
)

func main() {
	app := cli.App(Name, Descrip)

	ver := (Name + " " + Descrip + "\n" + runtime.GOOS + " " + runtime.GOARCH +
		"\n\n◇ Version: " + Version + "\n◇ Revision: " + CommitHash + "\n◇ Compiled: " + BuildTime)
	app.Version("v version", ver)

	dbg = *app.BoolOpt("d debug", false, "Enable debug logs")

	overseer.SetupLogBuilder(func(name string) overseer.Logger {
		return log.New(name)
	})
	logger.SetupLogBuilder(func(name string) logger.Logger {
		return log.New(name)
	})

	app.Command("list", "List all candidate files from the specified folder", cmdList)
	app.Command("one", "Generate code from a valid file and execute it", cmdRunOne)
	app.Command("up", "Convert all valid files from folder and execute them", cmdRunAll)

	app.Run(os.Args)
}

func cmdList(cmd *cli.Cmd) {
	cmd.Spec = "FOLDER"
	dir := cmd.StringArg("FOLDER", "", "the folder to list")

	cmd.Action = func() {
		files, err := parser.ParseFolder(*dir, false)
		if err != nil {
			fmt.Printf("List failed. Error: %v\n", err)
			return
		}

		for _, parsed := range files {
			if !parsed.IsValid() {
				if dbg {
					fmt.Printf("Invalid file: %s\n", parsed.Path)
				}
				continue
			}
			enabled := "●"
			if !parsed.Enabled {
				enabled = "■"
			}
			var langs []string
			for lng := range parsed.Blocks {
				langs = append(langs, lng)
			}
			fmt.Printf("%s %s ▻ %v\n", enabled, parsed.Path, langs)
		}
	}
}

func cmdRunOne(cmd *cli.Cmd) {
	cmd.Spec = "FILE [-f]"
	fname := cmd.StringArg("FILE", "", "the file to convert and run")
	force := cmd.BoolOpt("f force", false, "force by ignoring the header")

	cmd.Action = func() {
		fi, err := os.Stat(*fname)
		if err != nil {
			fmt.Printf("Run-one failed. Error: %v", err)
			return
		}
		if m := fi.Mode(); m.IsDir() || !m.IsRegular() || m&400 == 0 {
			fmt.Printf("The path must be a file. Path: %s", *fname)
			return
		}

		parseFile := parser.ParseFile(*fname)
		convFiles, err := parser.ConvertFile(parseFile, *force)
		if err != nil {
			fmt.Printf("Convert failed. Error: %v", err)
			return
		}
		if *force {
			fmt.Println("Unsafe mode enabled!")
		}

		ovr := overseer.NewOverseer()

		dir := filepath.Dir(*fname)
		baseLen := len(dir) + 1

		for lang, outFile := range convFiles {
			exe := parser.CodeBlocks[lang].Executable
			env := append(os.Environ(), "SPIN_FILE="+outFile)

			var p *overseer.Cmd
			if exe == "python3" {
				p = ovr.Add(parseFile.Id, exe, "-u", outFile[baseLen:])
			} else {
				p = ovr.Add(parseFile.Id, exe, outFile[baseLen:])
			}
			p.SetDir(dir)
			p.SetEnv(env)
			// TODO: maybe also DelayStart & RetryTimes?
		}

		fmt.Println("Starting proc. Press Ctrl+C to stop...")
		ovr.SuperviseAll()
		fmt.Println("\nShutdown.")
	}
}

func cmdRunAll(cmd *cli.Cmd) {
	cmd.Spec = "FOLDER [-n|--http] [--dry-run]"
	rootDir := cmd.StringArg("FOLDER", "", "the folder to convert and run")
	noHttp := cmd.BoolOpt("n no-http", false, "don't start the HTTP server")
	httpOpts := cmd.StringOpt("http", "localhost:12323", "HTTP server host:port")
	dryRun := cmd.BoolOpt("dry-run", false, "convert the folder and simulate running")

	cmd.Action = func() {
		dir := strings.TrimRight(*rootDir, "/")
		// This function will perform all folder checks
		pairs, err := parser.ConvertFolder(dir)
		if err != nil {
			fmt.Printf("Run-all failed. Error: %v", err)
			return
		}

		if *dryRun {
			*noHttp = true
		}
		ovr := overseer.NewOverseer()

		go func() {
			if *noHttp {
				fmt.Println("HTTP server disabled")
				return
			}
			// Setup HTTP server
			srv := http.NewServer(*httpOpts)
			// Enable Overseer endpoints
			http.OverseerEndpoint(srv, ovr)
			http.Serve(srv)
		}()

		baseLen := len(dir) + 1
		for infile, convFiles := range pairs {
			for lang, outFile := range convFiles {
				fmt.Printf("%s ==> %s\n", infile, outFile)

				exe := parser.CodeBlocks[lang].Executable
				env := append(os.Environ(), "SPIN_FILE="+outFile)
				p := ovr.Add(outFile, exe, outFile[baseLen:])
				p.SetDir(dir)
				p.SetEnv(env)
				// TODO: maybe also DelayStart & RetryTimes?
			}
		}

		if *dryRun {
			fmt.Println("Simulation over.")
		} else {
			fmt.Println("Starting procs. Press Ctrl+C to stop...")
			ovr.SuperviseAll()
			fmt.Println("\nShutdown.")
		}
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	. "github.com/ShinyTrinkets/meta-logger"
	"github.com/ShinyTrinkets/spinal/command"
	"github.com/ShinyTrinkets/spinal/parser"
	log "github.com/azer/logger"
	cli "github.com/jawher/mow.cli"
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

	SetupLogBuilder(func(name string) Logger {
		return log.New(name)
	})

	app.Command("list", "List all candidate source-files from folder", cmdList)
	app.Command("one", "Convert a source-file and execute it", cmdRunOne)
	app.Command("up", "Convert all source-files from folder and execute them", cmdRunAll)
	app.Command("check", "Show info about a running Spinal instance", cmdClient)

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

func cmdClient(cmd *cli.Cmd) {
	cmd.Spec = "[-c]"
	httpOpts := cmd.StringOpt("c http", "localhost:12323", "HTTP server host:port")

	cmd.Action = func() {
		resp, err := http.Get("http://" + *httpOpts + "/procs")
		if err != nil {
			fmt.Printf("Failed Spinal connection. Error: %v", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed Spinal reponse. Error: %v", err)
			return
		}
		fmt.Println("Procs: " + string(body))
	}
}

func cmdRunOne(cmd *cli.Cmd) {
	cmd.Spec = "FILE [-f]"
	fname := cmd.StringArg("FILE", "", "the file to convert and run")
	force := cmd.BoolOpt("f force", false, "force by ignoring the header")

	cmd.Action = func() {
		command.RunOne(*fname, *force)
	}
}

func cmdRunAll(cmd *cli.Cmd) {
	cmd.Spec = "FOLDER [-n|--http] [--dry-run]"
	rootDir := cmd.StringArg("FOLDER", "", "the folder to convert and run")
	noHTTP := cmd.BoolOpt("n no-http", false, "don't start the HTTP server")
	httpOpts := cmd.StringOpt("http", "localhost:12323", "HTTP server host:port")
	dryRun := cmd.BoolOpt("dry-run", false, "convert the folder and simulate running")

	cmd.Action = func() {
		command.RunAll(*rootDir, *httpOpts, *noHTTP, *dryRun)
	}
}

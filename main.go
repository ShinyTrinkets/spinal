package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/ShinyTrinkets/overseer.go"
	"github.com/ShinyTrinkets/spinal/http"
	"github.com/ShinyTrinkets/spinal/parser"
	"github.com/jawher/mow.cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func main() {
	app := cli.App(Name, Descrip)

	ver := (Name + " " + Descrip + "\n" + runtime.GOOS + " " + runtime.GOARCH +
		"\nVersion: " + Version + "\nRevision: " + CommitHash + "\nCompiled: " + BuildTime)
	app.Version("v version", ver)

	dbg := app.BoolOpt("d debug", false, "Enable debug logs")

	app.Before = func() {
		zerolog.TimeFieldFormat = ""
		zerolog.MessageFieldName = "m"
		// Pretty log enabled
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// Toggle debug logs
		if *dbg {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	}

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
			log.Fatal().Err(err).Msg("List failed")
		}
		log.Print("=== LIST ===")

		for _, parsed := range files {
			if !parsed.IsValid() {
				// log.Print("Invalid: " + parsed.Path)
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
			log.Info().Msgf("%s %s ▻ %v", enabled, parsed.Path, langs)
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
			log.Fatal().Err(err).Msg("Run-one failed")
		}
		if fi.IsDir() {
			log.Fatal().Msg("The path must be a file")
		}

		parseFile := parser.ParseFile(*fname)
		convFiles, err := parser.ConvertFile(parseFile, *force)
		if err != nil {
			log.Fatal().Err(err).Msg("Run-one failed")
		}
		log.Print("=== RUN-ONE ===")
		if *force {
			log.Warn().Msg("Unsafe mode enabled")
		}

		ovr := overseer.NewOverseer()

		dir := filepath.Dir(*fname)
		baseLen := len(dir) + 1

		for lang, outFile := range convFiles {
			exe := parser.CodeBlocks[lang].Executable
			p := ovr.Add(parseFile.Id, exe, outFile[baseLen:])
			p.SetDir(dir)
			// TODO: maybe also DelayStart & RetryTimes?
		}

		ovr.SuperviseAll()
	}
}

func cmdRunAll(cmd *cli.Cmd) {
	cmd.Spec = "FOLDER [-n|--http]"
	dir := cmd.StringArg("FOLDER", "", "the folder to convert and run")
	noHttp := cmd.BoolOpt("n no-http", false, "don't start the HTTP server")
	httpOpts := cmd.StringOpt("http", "localhost:12323", "HTTP server host:port")

	cmd.Action = func() {
		// This function will perform all folder checks
		pairs, err := parser.ConvertFolder(*dir)
		if err != nil {
			log.Fatal().Err(err).Msg("Run-all failed")
		}
		log.Print("=== RUN ALL ===")

		ovr := overseer.NewOverseer()
		go func() {
			if *noHttp {
				log.Info().Msg("HTTP server disabled")
				return
			}
			// Setup HTTP server
			srv := http.NewServer(*httpOpts)
			// Enable Overseer endpoints
			http.OverseerEndpoint(srv, ovr)
			http.Serve(srv)
		}()

		baseLen := len(*dir) + 1
		for infile, convFiles := range pairs {
			for lang, outFile := range convFiles {
				log.Info().Msgf("%s ==> %s", infile, outFile)

				exe := parser.CodeBlocks[lang].Executable
				env := append(os.Environ(), "SPIN_FILE="+outFile)
				p := ovr.Add(outFile, exe, outFile[baseLen:])
				p.SetDir(*dir)
				p.SetEnv(env)
				// TODO: maybe also DelayStart & RetryTimes?
			}
		}

		ovr.SuperviseAll()
	}
}

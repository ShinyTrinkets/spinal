package main

import (
	"os"
	"path"
	// "time"

	"github.com/ShinyTrinkets/gears.go/parser"
	"github.com/ShinyTrinkets/overseer.go"
	"github.com/jawher/mow.cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	Name    = "gears"
	Version = "0.0.1"
)

func main() {
	app := cli.App(Name, "Pretty dumb code runner")
	app.Version("v version", Name+" v"+Version)

	app.Spec = "[-d]"
	dbg := app.BoolOpt("d debug", false, "Debug mode enabled")

	app.Before = func() {
		zerolog.TimeFieldFormat = ""
		if *dbg {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	}

	app.Command("list", "List all candidate files from the specified folder", cmdList)
	app.Command("convert", "Convert all valid files from the specified folder", cmdConvert)
	app.Command("run-one", "Generate code from a valid file and execute it", cmdRunOne)

	app.Run(os.Args)
}

// Sample use: vault list OR vault config list
func cmdList(cmd *cli.Cmd) {
	cmd.Spec = "FOLDER"
	dir := cmd.StringArg("FOLDER", "", "the folder to list")

	cmd.Action = func() {
		files, err := parser.ParseFolder(*dir, false)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Print("=== LIST ===")

		for _, parsed := range files {
			if !parsed.IsValid() {
				// log.Print("Invalid: " + parsed.Path)
				continue
			}
			log.Info().Msgf("%s : %d lang", parsed.Path, len(parsed.Blocks))
		}
	}
}

func cmdConvert(cmd *cli.Cmd) {
	cmd.Spec = "FOLDER"
	dir := cmd.StringArg("FOLDER", "", "the folder to convert")

	cmd.Action = func() {
		// This function will perform all folder checks
		pairs, err := parser.ConvertFolder(*dir)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Print("=== CONVERT ===")

		baseLen := len(*dir) + 1
		for infile, outFiles := range pairs {
			for _, outFile := range outFiles {
				log.Info().Msgf("%s => %s\n", infile[baseLen:], outFile[baseLen:])
			}
		}
	}
}

func cmdRunOne(cmd *cli.Cmd) {
	cmd.Spec = "FILE"
	fname := cmd.StringArg("FILE", "", "the file to convert and run")

	cmd.Action = func() {
		fi, err := os.Stat(*fname)
		if err != nil {
			log.Fatal().Err(err)
		}
		if fi.IsDir() {
			log.Fatal().Msg("The path must be a file")
		}

		parseFile := parser.ParseFile(*fname)
		convFiles, err := parser.ConvertFile(parseFile)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Print("=== RUN-ONE ===")

		ovr := overseer.NewOverseer()

		dir := path.Dir(*fname)
		for lang, outFile := range convFiles {
			exe := parser.CodeBlocks[lang].Executable
			p := ovr.Add(parseFile.Id, exe, outFile)
			p.Dir = dir
		}

		ovr.SuperviseAll()
	}
}

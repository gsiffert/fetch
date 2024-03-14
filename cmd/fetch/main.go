// Package main define the configuration and the CLI for the fetch command.
package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	version = "0.1.0"
)

func main() {
	app := App{}

	cliApp := &cli.App{
		Name:    "fetch",
		Usage:   "Download the page of a website",
		Version: version,
		Before:  app.before,
		Action:  app.run,
		After:   app.after,
		Flags:   app.config.Flags(),
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

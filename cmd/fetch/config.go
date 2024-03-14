package main

import "github.com/urfave/cli/v2"

// Config holds the configuration for the CLI.
type Config struct {
	MetaData     bool
	DownloadPath string
	DSN          string
}

func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "metadata",
			Usage:       "fetch metadata for the given sites",
			Destination: &c.MetaData,
			Value:       false,
		},
		&cli.StringFlag{
			Name:        "dsn",
			Usage:       "DSN for the sqlite database",
			Destination: &c.DSN,
			Value:       "file:fetch.sqlite?cache=shared&mode=rwc",
			EnvVars:     []string{"FETCH_DSN"},
		},
		&cli.StringFlag{
			Name:        "download-path",
			Usage:       "Path to the directory to save the downloaded files",
			Destination: &c.DownloadPath,
			Value:       ".",
			EnvVars:     []string{"FETCH_DOWNLOAD_PATH"},
		},
	}
}

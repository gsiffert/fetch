package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gsiffert/fetch/internal/disk"
	"github.com/gsiffert/fetch/internal/fetcher"
	"github.com/gsiffert/fetch/internal/service"
	"github.com/gsiffert/fetch/internal/sqlite"
	"github.com/urfave/cli/v2"
)

// App implements the lifecycle of the CLI and forward the command to the service layer.
type App struct {
	config Config

	service      *service.Service
	metadataRepo *sqlite.MetaDataRepo
	logger       *slog.Logger
}

func (a *App) before(c *cli.Context) error {
	metadataRepo, err := sqlite.NewMetaDataRepo(c.Context, a.config.DSN)
	if err != nil {
		return fmt.Errorf("new metadata repo: %w", err)
	}

	a.metadataRepo = metadataRepo
	a.logger = slog.Default()
	f := fetcher.New(http.DefaultClient)
	d := disk.New(a.config.DownloadPath)
	a.service = service.New(f, d, a.logger, a.metadataRepo)

	return nil
}

func (a *App) metadataCommand(ctx context.Context, sites []string) error {
	metadataItems, err := a.service.GetMetaDataForSites(ctx, sites...)
	if err != nil {
		return fmt.Errorf("service get metadata for sites: %w", err)
	}

	var strs []string
	for _, metadata := range metadataItems {
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("site: %s\n", metadata.Site))
		builder.WriteString(fmt.Sprintf("num_links: %d\n", metadata.NumLinks))
		builder.WriteString(fmt.Sprintf("images: %d\n", metadata.NumImages))
		builder.WriteString(fmt.Sprintf("last_fetch: %s\n", metadata.LastFetched))
		strs = append(strs, builder.String())
	}
	fmt.Println(strings.Join(strs, "\n"))

	return nil
}

func (a *App) run(c *cli.Context) error {
	ctx := c.Context
	sites := c.Args().Slice()

	if a.config.MetaData {
		return a.metadataCommand(ctx, sites)
	}

	if err := a.service.Fetch(ctx, sites...); err != nil {
		return fmt.Errorf("service fetch: %w", err)
	}

	return nil
}

func (a *App) after(_ *cli.Context) error {
	if a.metadataRepo != nil {
		if err := a.metadataRepo.Close(); err != nil {
			return fmt.Errorf("close metadata repo: %w", err)
		}
	}
	return nil
}

// Package service implements the Service layer, connecting the infrastructure layer with the domain layer.
package service

import (
	"context"
	"io"
	"log/slog"

	"github.com/gsiffert/fetch/internal/domain"
)

// Disk defines the interface to save the content of a WebPage.
type Disk interface {
	NewPageWriter(ctx context.Context, name string) (io.WriteCloser, error)
}

// Fetcher defines the interface to download a WebPage.
type Fetcher interface {
	Fetch(ctx context.Context, site string) (*FetchedItem, error)
}

// MetaDataRepository defines the interface to save and retrieve domain.MetaData.
type MetaDataRepository interface {
	ByIDs(ctx context.Context, ids []domain.PageID) ([]domain.MetaData, error)
	Save(ctx context.Context, metaData domain.MetaData) error
}

// Service implements the functionality exposed to the application.
type Service struct {
	fetcher      Fetcher
	disk         Disk
	logger       *slog.Logger
	metaDataRepo MetaDataRepository
}

// New instantiate a new Service.
func New(fetcher Fetcher, disk Disk, logger *slog.Logger, metaDataRepo MetaDataRepository) *Service {
	return &Service{
		fetcher:      fetcher,
		disk:         disk,
		logger:       logger,
		metaDataRepo: metaDataRepo,
	}
}

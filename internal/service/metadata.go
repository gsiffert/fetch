package service

import (
	"context"
	"fmt"

	"github.com/gsiffert/fetch/internal/domain"
)

// GetMetaDataForSites retrieves a list of domain.MetaData from the given sites.
// It returns an error if it fails to retrieve the data from the repository.
func (s *Service) GetMetaDataForSites(ctx context.Context, sites ...string) ([]domain.MetaData, error) {
	ids := make([]domain.PageID, len(sites))
	for i, site := range sites {
		ids[i] = domain.PageID(site)
	}

	metadataItems, err := s.metaDataRepo.ByIDs(ctx, ids)
	if err != nil {
		s.logger.Error("Failed to get metadata.", "ids", ids, "error", err)
		return nil, fmt.Errorf("get metadata: %w", err)
	}

	return metadataItems, nil
}

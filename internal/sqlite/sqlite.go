// Package sqlite is part of the infrastructure layer, it implements the repositories using sqlite.
package sqlite

import (
	"context"
	"fmt"

	"github.com/gsiffert/fetch/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// MetaDataRepo implements the service.MetaDataRepository interface.
type MetaDataRepo struct {
	db *sqlx.DB
}

// NewMetaDataRepo instantiates a new MetaDataRepo.
func NewMetaDataRepo(ctx context.Context, dsn string) (*MetaDataRepo, error) {
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	repo := &MetaDataRepo{db: db}
	if err := repo.createTables(ctx); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return repo, nil
}

// Close the database connections.
func (r *MetaDataRepo) Close() error {
	return r.db.Close()
}

// This code will likely be removed in the future by using a migration tool.
func (r *MetaDataRepo) createTables(ctx context.Context) error {
	const query = `
	CREATE TABLE IF NOT EXISTS metadata (
	    id VARCHAR(255) PRIMARY KEY,
	    site VARCHAR(255) NOT NULL,
	    last_fetched DATETIME NOT NULL,
	    num_links INT UNSIGNED NOT NULL,
	    num_images INT UNSIGNED NOT NULL
	)
`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("create metadata table: %w", err)
	}

	return nil
}

// ByIDs retrieves a list od domain.MetaData matching the given ids.
func (r *MetaDataRepo) ByIDs(ctx context.Context, ids []domain.PageID) ([]domain.MetaData, error) {
	const baseQuery = `
	SELECT id, site, last_fetched, num_links, num_images
	FROM metadata
	WHERE id IN(?)
`

	query, args, err := sqlx.In(baseQuery, ids)
	if err != nil {
		return nil, fmt.Errorf("build sql in query: %w", err)
	}

	var items []domain.MetaData
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	for rows.Next() {
		var m domain.MetaData
		err = rows.Scan(&m.ID, &m.Site, &m.LastFetched, &m.NumLinks, &m.NumImages)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		items = append(items, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return items, nil
}

// Save the domain.MetaData.
func (r *MetaDataRepo) Save(ctx context.Context, m domain.MetaData) error {
	const query = `
	INSERT INTO metadata(id, site, last_fetched, num_links, num_images)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		last_fetched = ?,
		num_links = ?,
		num_images = ?
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		m.ID,
		m.Site,
		m.LastFetched,
		m.NumLinks,
		m.NumImages,
		m.LastFetched,
		m.NumLinks,
		m.NumImages,
	)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

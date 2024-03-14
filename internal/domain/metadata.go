// Package domain implements the domain layer, this package holds the domain logic and the domain models.
package domain

import "time"

// MetaData represents the metadata computed from a Page.
type MetaData struct {
	ID          PageID
	Site        string
	LastFetched time.Time
	NumLinks    int
	NumImages   int
}

package domain

import (
	"net/url"
	"path"
)

// PageID represents a unique ID for a Page.
type PageID string

// Strings implements the Stringer interface for a PageID.
func (id PageID) String() string {
	return string(id)
}

// Page represents a webpage.
type Page struct {
	ID           PageID
	Site         string
	FileLocation string
}

// NewPage instantiates a new Page.
func NewPage(u *url.URL) Page {
	site := path.Join(u.Host, u.Path)
	return Page{
		ID:           PageID(u.String()),
		Site:         site,
		FileLocation: url.PathEscape(site),
	}
}

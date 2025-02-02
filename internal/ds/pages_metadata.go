package ds

import (
	"github.com/a-h/templ"
)

type PagesMetadata []PageMetadata

type PageMetadata struct {
	Name string
	Path templ.SafeURL
}

func NewPageMetadata(name string, path string) PageMetadata {
	return PageMetadata{Name: name, Path: templ.URL(path)}
}

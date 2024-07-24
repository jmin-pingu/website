package ds

import (
	"github.com/a-h/templ"
)

type Pages []Page

type Page struct {
	Name string
	Link templ.SafeURL
}

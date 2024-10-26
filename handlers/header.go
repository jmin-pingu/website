package handlers

import (
	"mywebsite/ds"

	"github.com/a-h/templ"
)

var PAGES_METADATA ds.Pages = []ds.Page{
	ds.Page{Name: "Home", Link: templ.URL("/")},
	ds.Page{Name: "Blog", Link: templ.URL("/blog/")},
	ds.Page{Name: "Resources", Link: templ.URL("/resources/")},
	ds.Page{Name: "Projects", Link: templ.URL("/projects/")},
	ds.Page{Name: "Creative", Link: templ.URL("/creative/")},
}

var POSTS_METADATA ds.Posts = []*ds.Post{}

var POSTS_TAGS ds.OrderedList[string] = make(ds.OrderedList[string], 0)

var DISPLAY_TAGS ds.Set[string] = make(ds.Set[string], 0)

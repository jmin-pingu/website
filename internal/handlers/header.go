package handlers

import (
	"internal/db"
	"internal/ds"
)

var PAGES_METADATA ds.PagesMetadata = []ds.PageMetadata{
	ds.NewPageMetadata("Home", "/"),
	ds.NewPageMetadata("Blog", "/blog/"),
	ds.NewPageMetadata("TIL", "/til/"),
	ds.NewPageMetadata("Resources", "/resources/"),
	ds.NewPageMetadata("Projects", "/projects/"),
	ds.NewPageMetadata("Reading", "/reading_list/"),
	ds.NewPageMetadata("Creative", "/creative/"),
}

var BOOKS ds.StrictDict[string, db.Book]
var BOOK_TAGS ds.OrderedList[string] = make(ds.OrderedList[string], 0)

var DISPLAY_TAGS ds.Set[string] = make(ds.Set[string], 0)

var TIL_DISPLAY_TAGS ds.Set[string] = make(ds.Set[string], 0)

type BlogPostTags struct {
	Cs       string `json:"cs"`
	Personal string `json:"personal"`
}

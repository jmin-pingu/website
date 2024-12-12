package handlers

import (
	"mywebsite/db"
	"mywebsite/ds"
)

var PAGES_METADATA ds.PagesMetadata = []ds.PageMetadata{
	ds.NewPageMetadata("Home", "/"),
	ds.NewPageMetadata("Blog", "/blog/"),
	ds.NewPageMetadata("Resources", "/resources/"),
	//	ds.NewPageMetadata("Projects", "/projects/"),
	ds.NewPageMetadata("Creative", "/creative/"),
	ds.NewPageMetadata("Reading List", "/reading_list/"),
}

var POSTS_METADATA ds.PostsMetadata = []*ds.PostMetadata{}

var BOOKS ds.StrictDict[string, db.Book]
var BOOK_TAGS ds.OrderedList[string] = make(ds.OrderedList[string], 0)

var POSTS_TAGS ds.OrderedList[string] = make(ds.OrderedList[string], 0)

var DISPLAY_TAGS ds.Set[string] = make(ds.Set[string], 0)

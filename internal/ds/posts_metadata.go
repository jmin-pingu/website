package ds

import (
	"time"

	"github.com/a-h/templ"
)

type PostsMetadata []*PostMetadata

type PostMetadata struct {
	Title     string
	Date      time.Time
	PostID    int
	Path      templ.SafeURL
	Tags      Set[string]
	TagsFixed []string
	Display   bool
}

func NewPostMetadata(title string, date time.Time, path string, id int, tags Set[string]) *PostMetadata {
	return &PostMetadata{
		Title:     title,
		Date:      date,
		Path:      templ.URL(path),
		Tags:      tags,
		TagsFixed: tags.ToList(),
		PostID:    id,
		Display:   true,
	}
}

func (pm *PostsMetadata) AddPostMetadata(title string, date time.Time, id int, path string, tags Set[string]) {
	*pm = append(*pm, NewPostMetadata(title, date, path, id, tags))
}

func (pm *PostsMetadata) GetPostMetadata(path string) *PostMetadata {
	for _, post := range *pm {
		if post.Path == templ.SafeURL(path) {
			return post
		}
	}
	return nil
}

func (pm *PostsMetadata) ContainsPost(id int) bool {
	for _, post := range *pm {
		if post.PostID == id {
			return true
		}
	}
	return false
}

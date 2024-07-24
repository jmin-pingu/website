package ds

import (
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type Posts []*Post

type Post struct {
	Title   string
	Date    time.Time
	ID      uuid.UUID
	Link    templ.SafeURL
	Tags    Set[string]
	Display bool
}

func NewPost(title string, date time.Time, link string, tags Set[string]) *Post {
	return &Post{
		Title:   title,
		Date:    date,
		Link:    templ.URL(link),
		Tags:    tags,
		ID:      uuid.New(),
		Display: true,
	}
}

func (posts *Posts) AddPost(title string, date time.Time, link string, tags Set[string]) {
	*posts = append(*posts, NewPost(title, date, link, tags))
}

func (posts *Posts) GetPost(link string) *Post {
	for _, post := range *posts {
		if post.Link == templ.SafeURL(link) {
			return post
		}
	}
	return nil
}

func (posts *Posts) FilterPosts() {

}

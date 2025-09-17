package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"internal/db"
	"internal/ds"
	"internal/pub/pages"
	"internal/pub/pages/creative"
)

// NOTE: bare bones example of defining a handler, later if we want to keep track
// additional data per page
// type countHandler struct {
// 	mu sync.Mutex // guards n
// 	n  int
// }

// func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()
// 	h.n++
// 	fmt.Fprintf(w, "count is %d\n", h.n)
// }

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.HomePage(&PAGES_METADATA).Render(context.Background(), w)
}

type linksHandler struct{}

func (h *linksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.LinksPage(&PAGES_METADATA).Render(context.Background(), w)
}

type blogSearchHandler struct{}

func (h *blogSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: implement fuzzy searching https://github.com/lithammer/fuzzysearch
	if r.Method == http.MethodPost {
		log.Println("POST /blog/search: request received")
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		log.Printf("POST /blog/search: %s", r.Form["search"])
		pages.BlogPosts(&POSTS_METADATA, DISPLAY_TAGS, r.Form["search"][0]).Render(context.Background(), w)
	} else {
		pages.BlogPage(&PAGES_METADATA, &POSTS_METADATA, &POSTS_TAGS, DISPLAY_TAGS).Render(context.Background(), w)
	}
}

type blogHandler struct{}

func (h *blogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("POST /blog/: request received")

		// NOTE: go/http doesn't guarantee you can write a response before reading the
		// request.  Some servers will do it and some clients will do it, but not all, so we
		// intentionally disallow it, to prevent surprises.  It's easier for us to relax this
		// restriction in the future if we discover we're wrong (spec references and comprehensive
		// testing of all major clients & servers welcome!) than it is for us to allow it now and
		// then take it away.

		// References: reference ParseForm() and then r.Form and r.PostForm in go/http
		// - https://stackoverflow.com/questions/55945248/go-keep-get-scan-file-error-http-invalid-read-on-closed-body
		// - https://pkg.go.dev/net/http#Request.FormValue
		// - https://github.com/golang/go/issues/4637
		// - "ParseForm populates r.Form and r.PostForm. For all requests, ParseForm parses the raw query from the URL and updates r.Form."
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		log.Println(r.Form)

		for k, v := range r.Form {
			log.Printf("POST /blog/: key %v, value %v\n", k, v)
			if v[0] == "1" {
				DISPLAY_TAGS.Add(k)
			} else {
				DISPLAY_TAGS.Remove(k)
			}
		}
		pages.BlogPosts(&POSTS_METADATA, DISPLAY_TAGS, "").Render(context.Background(), w)
	} else {
		RenderPosts()
		pages.BlogPage(&PAGES_METADATA, &POSTS_METADATA, &POSTS_TAGS, DISPLAY_TAGS).Render(context.Background(), w)
	}
}

type resourcesHandler struct{}

func (h *resourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.ResourcesPage(&PAGES_METADATA).Render(context.Background(), w)
}

type projectsHandler struct{}

func (h *projectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.ProjectsPage(&PAGES_METADATA).Render(context.Background(), w)
}

type creativeHandler struct{}

func (h *creativeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	creative.CreativePage(&PAGES_METADATA).Render(context.Background(), w)
}

type readingListHandler struct{}

func (h *readingListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: add cache-ing
	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to website: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	BOOKS = db.GetBooks(dbpool)
	pages.ReadingListPage(&PAGES_METADATA, &BOOKS).Render(context.Background(), w)
}

func SetUpRoutes() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("internal/pub/assets"))))
	http.Handle("/", new(homeHandler))
	http.Handle("/blog/", new(blogHandler))
	http.Handle("/blog/search", new(blogSearchHandler))
	http.Handle("/resources/", new(resourcesHandler))
	http.Handle("/projects/", new(projectsHandler))
	http.Handle("/creative/", new(creativeHandler))
	http.Handle("/reading_list/", new(readingListHandler))
	http.Handle("/links/", new(linksHandler))

}

func RenderPosts() {
	dbpool, err := db.GetConnection(os.Getenv("POSTGRES_DB"))
	if err != nil {
		log.Printf("db: %v", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	posts := db.GetPosts(dbpool)
	for _, post := range posts {
		tags := make(ds.Set[string], 0)
		for _, tag := range post.Tags {
			t := strings.ToLower(strings.TrimSpace(tag))
			POSTS_TAGS.Add(t) // Universal tag set: ds.OrderedList
			tags.Add(t)       // Per post tag set: ds.Set
		}

		if !POSTS_METADATA.ContainsPost(post.PostID) {
			url := "/blog/" + post.Link
			POSTS_METADATA.AddPostMetadata(post.Title, post.Date, post.PostID, url, tags)
			http.HandleFunc(
				url,
				func(w http.ResponseWriter, _ *http.Request) {
					pages.BlogEntryPage(post.Title, post.Content, &PAGES_METADATA, &POSTS_METADATA).Render(context.Background(), w)
				})
		}
	}
}

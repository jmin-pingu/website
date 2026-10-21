package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"internal/db"
	"internal/ds"
	"internal/pub/pages"
	"internal/pub/pages/creative"

	"github.com/jackc/pgx/v5/pgxpool"
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

// DBPOOL is opened once in cmd/run and shared by every request
var DBPOOL *pgxpool.Pool

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pages.HomePage(&PAGES_METADATA).Render(context.Background(), w)
}

type blogSearchHandler struct{}

func (h *blogSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: implement fuzzy searching https://github.com/lithammer/fuzzysearch
	posts, tags, err := loadPosts()
	if err != nil {
		serverError(w, err)
		return
	}
	if r.Method == http.MethodPost {
		log.Println("POST /blog/search: request received")
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		log.Printf("POST /blog/search: %s", r.Form["search"])
		pages.BlogPosts(&posts, DISPLAY_TAGS, r.Form.Get("search")).Render(context.Background(), w)
	} else {
		pages.BlogPage(&PAGES_METADATA, &posts, &tags, DISPLAY_TAGS).Render(context.Background(), w)
	}
}

type blogHandler struct{}

// Serves the post list at /blog/ and each post at /blog/<link>. Posts are read
// from the db on every request, so uploads, edits and deletes show up without
// a restart.
func (h *blogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog/" {
		post, err := db.GetPostByLink(DBPOOL, strings.TrimPrefix(r.URL.Path, "/blog/"))
		if err != nil {
			serverError(w, err)
			return
		}
		if post == nil {
			http.NotFound(w, r)
			return
		}
		pages.BlogEntryPage(post.Title, post.Content, &PAGES_METADATA, nil).Render(context.Background(), w)
		return
	}

	posts, tags, err := loadPosts()
	if err != nil {
		serverError(w, err)
		return
	}
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
		pages.BlogPosts(&posts, DISPLAY_TAGS, "").Render(context.Background(), w)
	} else {
		pages.BlogPage(&PAGES_METADATA, &posts, &tags, DISPLAY_TAGS).Render(context.Background(), w)
	}
}

type tilSearchHandler struct{}

func (h *tilSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entries, tags, err := loadTIL()
	if err != nil {
		serverError(w, err)
		return
	}
	if r.Method == http.MethodPost {
		log.Println("POST /til/search: request received")
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		pages.TILPosts(&entries, TIL_DISPLAY_TAGS, r.Form.Get("search")).Render(context.Background(), w)
	} else {
		pages.TILPage(&PAGES_METADATA, &entries, &tags, TIL_DISPLAY_TAGS).Render(context.Background(), w)
	}
}

type tilHandler struct{}

// Same as blogHandler: the list at /til/ and each entry at /til/<link>
func (h *tilHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/til/" {
		entry, err := db.GetTILByLink(DBPOOL, strings.TrimPrefix(r.URL.Path, "/til/"))
		if err != nil {
			serverError(w, err)
			return
		}
		if entry == nil {
			http.NotFound(w, r)
			return
		}
		pages.BlogEntryPage(entry.Title, entry.Content, &PAGES_METADATA, nil).Render(context.Background(), w)
		return
	}

	entries, tags, err := loadTIL()
	if err != nil {
		serverError(w, err)
		return
	}
	if r.Method == http.MethodPost {
		log.Println("POST /til/: request received")

		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		for k, v := range r.Form {
			if v[0] == "1" {
				TIL_DISPLAY_TAGS.Add(k)
			} else {
				TIL_DISPLAY_TAGS.Remove(k)
			}
		}
		pages.TILPosts(&entries, TIL_DISPLAY_TAGS, "").Render(context.Background(), w)
	} else {
		pages.TILPage(&PAGES_METADATA, &entries, &tags, TIL_DISPLAY_TAGS).Render(context.Background(), w)
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
	BOOKS = db.GetBooks(DBPOOL)
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
	http.Handle("/til/", new(tilHandler))
	http.Handle("/til/search", new(tilSearchHandler))
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("db: %v", err)
	http.Error(w, "something went wrong", http.StatusInternalServerError)
}

// tagSet normalizes a post's tags and adds each one to the list of all tags
func tagSet(raw []string, all *ds.OrderedList[string]) ds.Set[string] {
	tags := make(ds.Set[string], 0)
	for _, tag := range raw {
		t := strings.ToLower(strings.TrimSpace(tag))
		all.Add(t)  // Universal tag set: ds.OrderedList
		tags.Add(t) // Per post tag set: ds.Set
	}
	return tags
}

// loadPosts builds the post list and tag list from the db
func loadPosts() (ds.PostsMetadata, ds.OrderedList[string], error) {
	metadata := ds.PostsMetadata{}
	all_tags := ds.OrderedList[string]{}
	posts, err := db.GetPosts(DBPOOL)
	if err != nil {
		return metadata, all_tags, err
	}
	for _, post := range posts {
		tags := tagSet(post.Tags, &all_tags)
		metadata.AddPostMetadata(post.Title, post.Date, db.Hash(post.PostID).String(), "/blog/"+post.Link, tags)
	}
	return metadata, all_tags, nil
}

// loadTIL builds the TIL list and tag list from the db
func loadTIL() (ds.PostsMetadata, ds.OrderedList[string], error) {
	metadata := ds.PostsMetadata{}
	all_tags := ds.OrderedList[string]{}
	entries, err := db.GetTIL(DBPOOL)
	if err != nil {
		return metadata, all_tags, err
	}
	for _, entry := range entries {
		tags := tagSet(entry.Tags, &all_tags)
		metadata.AddPostMetadata(entry.Title, entry.Date, db.Hash(entry.TILID).String(), "/til/"+entry.Link, tags)
	}
	return metadata, all_tags, nil
}

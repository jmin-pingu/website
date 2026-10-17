MAKEFLAGS += --silent

DB_USER ?= jonathanmin
DB_NAME ?= websitedb
# Skip empty drafts; the upload tool panics on posts without front matter
POSTS ?= $(shell find data/blog_posts -name '*.md' -size +0)

# Local targets use the Postgres on this machine instead of the remote one in
# .env. godotenv doesn't override variables that are already set.
db-setup db-load db-load-posts dev: export POSTGRES_IP := localhost
db-setup db-load db-load-posts dev: export POSTGRES_PORT := 5432

run:
	templ generate
	npm run tailwind-build
	go run cmd/run/main.go

# Create the local user and database if they don't exist yet
db-setup:
	psql -h localhost -d postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$(DB_USER)'" | grep -q 1 \
		|| createuser -h localhost $(DB_USER)
	psql -h localhost -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 \
		|| createdb -h localhost -O $(DB_USER) $(DB_NAME)

# Load posts and books into the local database (POSTS=path/to/post.md for one)
db-load: db-setup db-load-posts
	go run ./cmd/upload books data/books.json

db-load-posts:
	go run ./cmd/upload posts $(POSTS)

# Watch templ, tailwind, go and posts against the local database; open localhost:3000
dev: db-load
	npm run dev

# Upload a post to the remote database from .env, e.g. make publish POST=data/blog_posts/x.md
publish:
	test -n "$(POST)" || (echo "usage: make publish POST=data/blog_posts/<post>.md" && exit 1)
	go run ./cmd/upload posts $(POST)

build:
	templ generate
	npm run tailwind-build
	go build -o ./bin/run ./cmd/run
	go build -o ./bin/upload ./cmd/upload
	docker build -t ghcr.io/jmin-pingu/website/server:latest .
	docker push ghcr.io/jmin-pingu/website/server:latest

watch:
	npm run tailwind

push-image: build
	docker push ghcr.io/jmin-pingu/website/server:latest


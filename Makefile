MAKEFLAGS += --silent

DB_USER ?= jonathanmin
DB_NAME ?= websitedb

# Local targets use the Postgres on this machine instead of the remote one in
# .env. godotenv doesn't override variables that are already set.
db-setup db-load db-sync dev: export POSTGRES_IP := localhost
db-setup db-load db-sync dev: export POSTGRES_PORT := 5432

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

# Load posts, TILs and books into the local database
db-load: db-setup db-sync
	go run ./cmd/upload books data/books.json

# Make the local posts and til tables match data/posts and data/til
db-sync:
	go run ./cmd/upload sync data

# Watch templ, tailwind, go and posts against the local database; open localhost:3000
dev: db-load
	npm run dev

# Make the remote database from .env match data/: uploads new and edited
# posts and TILs, and deletes ones whose file is gone. No restart needed.
sync:
	go run ./cmd/upload sync data

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


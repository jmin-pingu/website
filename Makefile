MAKEFLAGS += --silent

run:
	templ generate
	npm run tailwind-build
	go run cmd/run/main.go

dev:
	npm run dev

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


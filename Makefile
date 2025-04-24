MAKEFLAGS += --silent

run:
	templ generate 
	npm run tailwind-build
	# open http://localhost:1323/
	go run cmd/run/main.go

build:
	templ generate 
	npm run tailwind-build
	go build -o ./bin/run ./cmd/run 
	go build -o ./bin/upload ./cmd/upload 
	docker build -t ghcr.io/jmin-pingu/website/server:latest .

watch:
	npm run tailwind

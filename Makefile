MAKEFLAGS += --silent

run:
	templ generate 
	# open http://localhost:1323/
	go run cmd/run/main.go

build:
	templ generate 
	go build ./cmd/run
	go build ./cmd/upload
	mv run bin
	mv upload bin

watch:
	npm run tailwind

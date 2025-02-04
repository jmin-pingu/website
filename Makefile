MAKEFLAGS += --silent

run:
	templ generate 
	# open http://localhost:1323/
	go run cmd/run/main.go

build:
	templ generate 
	go build -o ./bin/run ./cmd/run 
	go build -o ./bin/upload ./cmd/upload 

watch:
	npm run tailwind

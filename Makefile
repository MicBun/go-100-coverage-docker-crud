start:
	go run ./bin/web

build:
	go build ./bin/web

swag-init:
	swag init -g ./bin/web/main.go

test:
	go test -v ./...

coverage:
	go test -coverprofile='coverage.out' ./...
	go tool cover -html='coverage.out'
	del coverage.out

docker:
	docker build -t web .

docker-run:
	docker run -p 8080:8080 web

docker-compose:
	docker-compose up
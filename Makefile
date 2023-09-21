.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/go-tgbot ./cmd/app/main.go && cp ./config.yml ./.bin/config.yml

run: build
	docker-compose up --remove-orphans go-tgbot

rebuild:
	docker-compose up -d --no-deps --build
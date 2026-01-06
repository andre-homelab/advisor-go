SWAG ?= swag

.PHONY: swagger migrate api dev

swagger:
	$(SWAG) init -g inputs/api/main.go -o docs || \
		go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g inputs/api/main.go -o docs

migrate:
	go run ./cmd/migrate

api:
	go run . api

dev: swagger migrate api

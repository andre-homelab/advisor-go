SWAG ?= swag

.PHONY: docs migrate api dev

docs:
	$(SWAG) init -g inputs/api/main.go -o docs || \
		go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g inputs/api/main.go -o docs

api:
	go run . api

compose:
	docker compose up -d

dev: compose docs api 

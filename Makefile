.PHONY: run test

run:
	@go run main.go

test:
	@go test -v ./...

compose-up:
	@docker compose up --build -d --remove-orphans

compose-stop:
	@docker compose stop

compose-down:
	@docker compose down --remove-orphans

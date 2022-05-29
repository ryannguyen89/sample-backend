##@ Testing

.PHONY: test
test: ## Run unit test
	go test ./... -cover -race

test-cover: ## Run unit test with report coverage
	@go test -race -coverpkg=./... $(shell go list ./...) -coverprofile=cover.out
	@go tool cover -func=cover.out
	@rm -rf cover.out

##@ Developement

run: ## Run app locally
	go run main.go

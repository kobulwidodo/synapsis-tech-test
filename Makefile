.PHONY: run
run:
	@go run ./src/cmd/main.go

.PHONY: run-tests
run-tests:
	@go clean -cache
	@go test -v -failfast `go list ./... | grep -i 'business'` -cover

.PHONY: swaggo
swaggo:
	@/bin/rm -rf ./docs/swagger
	@`go env GOPATH`/bin/swag init -g ./src/cmd/main.go -o ./docs/swagger --parseInternal

.PHONY: swag-install
swag-install:
	@go install github.com/swaggo/swag/cmd/swag@v1.6.7

.PHONE: run-app
run-app:
	@make swaggo
	@make run

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source src/business/domain/$(domain)/$(domain).go -destination src/business/domain/mock/$(domain)/$(domain).go

.PHONY: mock-lib
mock-lib:
	@`go env GOPATH`/bin/mockgen -source src/lib/$(domain)/$(domain).go -destination src/lib/tests/mock/$(domain)/$(domain).go

.PHONY: mock-all
mock-all:
	@make mock-lib domain=auth
	@make mock domain=product
	@make mock domain=category
	@make mock domain=cart
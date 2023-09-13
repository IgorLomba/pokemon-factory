SRC=$(shell find . -name "*.go")
IMAGE_NAME=pokemon-factory
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1
linter:
	golangci-lint run -v

check-dependencies:
	go mod tidy
	git diff --exit-code go.mod
	git diff --exit-code go.sum

checks: test check-dependencies linter

up: build
	docker compose up -d

down:
	docker compose down

clean:
	docker compose down --volumes --remove-orphans --rmi all
	docker compose rm --force --stop --volumes

test:
	go test -cover ./...

deps: check-dependencies
	go mod tidy
	go mod download

run:
	go run main.go

# https://goswagger.io/install.html
gen-swagger:
	swagger generate spec -m -o swagger.json

serve-swagger:
	swagger serve -F swagger swagger.json

gen-serve-swagger:
	swagger generate spec -m -o swagger.json && swagger serve -F swagger swagger.json

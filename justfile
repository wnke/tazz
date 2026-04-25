set shell := ["bash", "-cu"]

gen-api:
	oapi-codegen -config internal/driver/http/oapi-codegen.yaml internal/driver/http/openapi.yaml

test:
	go test ./...

lint:
	golangci-lint run

check:
	just lint
	just test

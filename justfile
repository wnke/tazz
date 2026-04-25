set shell := ["bash", "-cu"]

generator := "oapi-codegen"

gen-api:
	{{generator}} -config internal/driver/http/oapi-codegen.yaml internal/driver/http/openapi.yaml

test:
	go test ./...

.PHONY: run swag-generate all

all: swag-generate run

run:
	go run cmd/main.go

swag-generate:
	cd cmd && swag init -g ../cmd/main.go -d ../config,../internal/models,../internal/controllers,../internal/storage/database,../internal/storage/repository -o ../docs

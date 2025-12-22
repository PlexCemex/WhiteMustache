.PHONY: build push help docker-build docker-push build-and-push

help:
	@echo "Available commands:"
	@echo "  make build        - Build Linux binary"
	@echo "  make docker-build - Build Docker image"
	@echo "  make push         - Push Docker image"
	@echo "  make docker-push  - Build and push Docker image"

build:
	cd ./mock-server && powershell -Command "$$env:CGO_ENABLED='0'; $$env:GOOS='linux'; $$env:GOARCH='amd64'; go build -o main main.go"

docker-build: build
	docker build -t plexcemex/whitemustache:latest .
	@echo "Docker image built"

push:
	docker push plexcemex/whitemustache:latest

docker-push: docker-build push
	@echo "Docker image pushed"

build-and-push: docker-push
	@echo "Build and push completed"

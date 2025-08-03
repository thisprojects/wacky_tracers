# Build variables
VERSION ?= dev
REGISTRY ?= devmesh
PLATFORMS ?= linux/amd64,linux/arm64

# Go variables
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

.PHONY: help build test clean docker deploy

help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## Development
build: ## Build all binaries
	@echo "Building DevMesh binaries..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/tracer cmd/tracer/main.go
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/collector cmd/collector/main.go
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/webhook cmd/webhook/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -tags=integration ./test/integration/...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out coverage.html

## Docker
docker-build: ## Build Docker images
	@echo "Building Docker images..."
	@docker build -f Dockerfile.tracer -t $(REGISTRY)/devmesh-tracer:$(VERSION) .
	@docker build -f Dockerfile.collector -t $(REGISTRY)/devmesh-collector:$(VERSION) .
	@docker build -f Dockerfile.webhook -t $(REGISTRY)/devmesh-webhook:$(VERSION) .

docker-push: ## Push Docker images
	@echo "Pushing Docker images..."
	@docker push $(REGISTRY)/devmesh-tracer:$(VERSION)
	@docker push $(REGISTRY)/devmesh-collector:$(VERSION)
	@docker push $(REGISTRY)/devmesh-webhook:$(VERSION)

## Kubernetes
generate: ## Generate Kubernetes manifests
	@echo "Generating manifests..."
	@controller-gen rbac:roleName=devmesh-controller crd webhook paths="./..." output:crd:artifacts:config=deployments/k8s/crds

deploy-dev: ## Deploy to development cluster
	@echo "Deploying to development cluster..."
	@kubectl apply -f deployments/k8s/
	@helm upgrade --install devmesh deployments/helm/ -n devmesh-system --create-namespace

## Protobuf (if using gRPC)
proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go-grpc_out=. api/proto/trace.proto

## Release
release: clean test build docker-build ## Build release artifacts
	@echo "Release $(VERSION) ready"
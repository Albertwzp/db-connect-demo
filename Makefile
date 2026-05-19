all: build

tidy:
	@echo "tidying modules..."
	go mod tidy

build:
	@echo "building single binary..."
	go build -o db-connect-demo

# Cross-compile for Linux amd64 (requires CGO_ENABLED=0 for cross-compilation without a Linux cross-compiler toolchain)
# NOTE: CGO_ENABLED=0 disables sqlite3 support (mattn/go-sqlite3 requires CGO).
#       To build with full sqlite3 support for Linux, set CGO_ENABLED=1 and provide
#       a cross-compiler via CC=x86_64-linux-gnu-gcc, or build natively on Linux.
build-linux:
	@echo "cross-compiling for linux/amd64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags="-w -s" -o db-connect-demo-linux
	@echo "✓ Built db-connect-demo-linux"

clean:
	@rm -f db-connect-demo
	@rm -f db-connect-demo-linux
	@rm -rf frontend/dist
	@echo "✓ Cleaned"

# Clean extra files created during development
cleanup:
	@echo "Removing development files..."
	@rm -f operator_main.go init_dirs.go test_operator.sh deploy.sh validate.go cleanup.go backends.json db-bench.exe README_NEW.md
	@echo "✓ Cleanup complete"

ui-build:
	@if [ -d frontend ]; then \
		echo "building frontend..."; \
		cd frontend && \
		if [ ! -d node_modules ]; then npm install; fi && \
		npm run build && \
		cd ..; \
	else \
		echo "frontend dir not found"; \
	fi

ui-clean:
	@rm -rf frontend/dist
	@echo "✓ Frontend cleaned"

# Standalone mode (no K8s required)
run: build ui-build
	@echo "Starting combined Operator+API server on port 8080"
	@echo "UI: http://localhost:8080/ui"
	@echo "API: http://localhost:8080/ping"
	./db-connect-demo -port=8080

# Kubernetes deployment
k8s-install:
	@echo "Installing CRD..."
	kubectl apply -f config/crd/config_crd_crds.yaml
	@echo "✓ CRDs installed"

k8s-rbac:
	@echo "Installing RBAC..."
	kubectl apply -f config/rbac/config_rbac_rbac.yaml
	@echo "✓ RBAC configured"

k8s-deploy: build docker-build k8s-install k8s-rbac
	@echo "Deploying Operator and API Server..."
	kubectl apply -f config/manager/config_manager_deployment.yaml
	@echo "✓ Deployed"

k8s-samples:
	@echo "Creating sample connections..."
	kubectl apply -f config/samples/config_samples_connections.yaml
	@echo "✓ Samples created"

k8s-uninstall:
	@echo "Uninstalling..."
	kubectl delete -f config/manager/config_manager_deployment.yaml --ignore-not-found
	kubectl delete -f config/rbac/config_rbac_rbac.yaml --ignore-not-found
	kubectl delete -f config/crd/config_crd_crds.yaml --ignore-not-found
	@echo "✓ Uninstalled"

# Docker
docker-build: build
	@echo "Building Docker image..."
	docker build -t db-connect-demo:latest .

docker-push:
	@echo "Pushing Docker image..."
	docker push db-connect-demo:latest

help:
	@echo "db-connect-demo - Combined Operator + API Server"
	@echo ""
	@echo "Standalone mode (no K8s required):"
	@echo "  make run       - Build and run combined service"
	@echo "  make build     - Build binary only"
	@echo "  make build-linux - Cross-compile for linux/amd64"
	@echo ""
	@echo "Kubernetes mode:"
	@echo "  make k8s-deploy     - Deploy to K8s cluster"
	@echo "  make k8s-samples    - Create sample connections"
	@echo "  make k8s-uninstall  - Remove from K8s"
	@echo ""
	@echo "Build:"
	@echo "  make ui-build  - Build React frontend"
	@echo "  make tidy      - Sync Go dependencies"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make cleanup   - Remove development files"
	@echo "  make ui-clean  - Remove frontend dist"
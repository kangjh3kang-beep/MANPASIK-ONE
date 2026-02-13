# ManPaSik 프로젝트 Makefile
# 빌드, 테스트, 배포 자동화

.PHONY: all build test lint clean docker-build docker-push k8s-apply help

# 변수
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Go 설정
GO := go
GOFLAGS := -v

# Docker 설정
DOCKER_REGISTRY ?= docker.io/manpasik
DOCKER_TAG ?= $(VERSION)

# Rust 설정
CARGO := cargo

# 서비스 목록
GO_SERVICES := gateway auth-service measurement-service user-service device-service
RUST_CRATES := manpasik-engine flutter-bridge

#==============================================================================
# 기본 타겟
#==============================================================================

all: build

help:
	@echo "ManPaSik 프로젝트 Makefile"
	@echo ""
	@echo "사용법:"
	@echo "  make build          - 모든 서비스 빌드"
	@echo "  make build-go       - Go 서비스 빌드"
	@echo "  make build-rust     - Rust 코어 빌드"
	@echo "  make test           - 모든 테스트 실행"
	@echo "  make test-go        - Go 테스트"
	@echo "  make test-rust      - Rust 테스트"
	@echo "  make lint           - 린트 검사"
	@echo "  make docker-build   - Docker 이미지 빌드"
	@echo "  make docker-push    - Docker 이미지 푸시"
	@echo "  make k8s-apply      - Kubernetes 배포"
	@echo "  make proto          - gRPC Proto 컴파일"
	@echo "  make clean          - 빌드 결과물 삭제"
	@echo "  make dev            - 개발 환경 시작"
	@echo ""

#==============================================================================
# 빌드
#==============================================================================

build: build-rust build-go
	@echo "✅ 전체 빌드 완료"

build-go:
	@echo "🔨 Go 서비스 빌드..."
	cd backend && $(GO) build $(GOFLAGS) $(LDFLAGS) -o ../bin/gateway ./gateway/cmd
	cd backend && $(GO) build $(GOFLAGS) $(LDFLAGS) -o ../bin/auth-service ./services/auth-service/cmd
	cd backend && $(GO) build $(GOFLAGS) $(LDFLAGS) -o ../bin/measurement-service ./services/measurement-service/cmd
	@echo "✅ Go 빌드 완료"

build-rust:
	@echo "🦀 Rust 코어 빌드..."
	cd rust-core && $(CARGO) build --release
	@echo "✅ Rust 빌드 완료"

#==============================================================================
# 테스트
#==============================================================================

test: test-rust test-go
	@echo "✅ 전체 테스트 완료"

test-go:
	@echo "🧪 Go 테스트..."
	cd backend && $(GO) test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Go 테스트 완료"

test-rust:
	@echo "🧪 Rust 테스트..."
	cd rust-core && $(CARGO) test --all
	@echo "✅ Rust 테스트 완료"

test-integration:
	@echo "🧪 통합 테스트 (backend/tests/e2e)..."
	cd backend && $(GO) test -v -tags=integration ./tests/e2e/...
	@echo "✅ 통합 테스트 완료"

#==============================================================================
# 린트
#==============================================================================

lint: lint-go lint-rust
	@echo "✅ 린트 완료"

lint-go:
	@echo "🔍 Go 린트..."
	cd backend && golangci-lint run ./...

lint-rust:
	@echo "🔍 Rust 린트..."
	cd rust-core && $(CARGO) clippy --all-targets -- -D warnings
	cd rust-core && $(CARGO) fmt --all -- --check

#==============================================================================
# Proto 컴파일
#==============================================================================

# Google well-known types. 기본 /usr/include (Linux). 실패 시 PROTO_GOOGLE_INCLUDE 설정.
PROTO_GOOGLE_INCLUDE ?= /usr/include

proto:
	@echo "📝 gRPC Proto 컴파일..."
	@command -v protoc >/dev/null 2>&1 || { echo "❌ protoc 없음. 설치 후: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; exit 1; }
	cd backend && PATH="$$PATH:$$(go env GOPATH)/bin" protoc \
		--proto_path=shared/proto \
		--proto_path=$(PROTO_GOOGLE_INCLUDE) \
		--go_out=. --go_opt=module=github.com/manpasik/backend \
		--go-grpc_out=. --go-grpc_opt=module=github.com/manpasik/backend \
		shared/proto/manpasik.proto shared/proto/health.proto
	@echo "✅ Proto 컴파일 완료 (E2E TestMeasurementFlow는 이 생성 코드 필요)"

#==============================================================================
# Docker
#==============================================================================

docker-build:
	@echo "🐳 Docker 이미지 빌드..."
	docker build -t $(DOCKER_REGISTRY)/gateway:$(DOCKER_TAG) -f backend/gateway/Dockerfile .
	docker build -t $(DOCKER_REGISTRY)/auth-service:$(DOCKER_TAG) -f backend/services/auth-service/Dockerfile .
	docker build -t $(DOCKER_REGISTRY)/measurement-service:$(DOCKER_TAG) -f backend/services/measurement-service/Dockerfile .
	@echo "✅ Docker 빌드 완료"

docker-push:
	@echo "📤 Docker 이미지 푸시..."
	docker push $(DOCKER_REGISTRY)/gateway:$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/auth-service:$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/measurement-service:$(DOCKER_TAG)
	@echo "✅ Docker 푸시 완료"

#==============================================================================
# Kubernetes
#==============================================================================

k8s-apply:
	@echo "☸️ Kubernetes 배포..."
	kubectl apply -f infrastructure/kubernetes/base/namespace.yaml
	kubectl apply -f infrastructure/kubernetes/base/config/
	kubectl apply -f infrastructure/kubernetes/base/services/
	@echo "✅ Kubernetes 배포 완료"

k8s-delete:
	@echo "🗑️ Kubernetes 리소스 삭제..."
	kubectl delete -f infrastructure/kubernetes/base/services/ --ignore-not-found
	kubectl delete -f infrastructure/kubernetes/base/config/ --ignore-not-found
	@echo "✅ Kubernetes 삭제 완료"

#==============================================================================
# 개발 환경
#==============================================================================

# Docker Compose: V2(docker compose) 기본. V1만 있으면 make DOCKER_COMPOSE=docker-compose make dev
DOCKER_COMPOSE ?= docker compose

dev:
	@echo "🚀 개발 환경 시작..."
	cd infrastructure/docker && $(DOCKER_COMPOSE) -f docker-compose.dev.yml up -d
	@echo "✅ 개발 환경 시작 완료"
	@echo "서비스 상태: $(DOCKER_COMPOSE) -f infrastructure/docker/docker-compose.dev.yml ps"

dev-stop:
	@echo "🛑 개발 환경 중지..."
	cd infrastructure/docker && $(DOCKER_COMPOSE) -f docker-compose.dev.yml down
	@echo "✅ 개발 환경 중지 완료"

dev-logs:
	cd infrastructure/docker && $(DOCKER_COMPOSE) -f docker-compose.dev.yml logs -f

#==============================================================================
# 정리
#==============================================================================

clean:
	@echo "🧹 정리..."
	rm -rf bin/
	rm -rf backend/coverage.out
	cd rust-core && $(CARGO) clean
	@echo "✅ 정리 완료"

#==============================================================================
# 유틸리티
#==============================================================================

deps:
	@echo "📦 의존성 설치..."
	cd backend && $(GO) mod download
	cd rust-core && $(CARGO) fetch
	@echo "✅ 의존성 설치 완료"

fmt:
	@echo "✨ 코드 포맷팅..."
	cd backend && $(GO) fmt ./...
	cd rust-core && $(CARGO) fmt --all
	@echo "✅ 포맷팅 완료"

version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

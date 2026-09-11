BINARY := iprbooks-dumper
CMD     := ./cmd/dumper
IMAGE   := iprbooks-dumper:latest

GOBIN         := $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(GOBIN)/golangci-lint
GO_ARCH_LINT  := $(GOBIN)/go-arch-lint

.PHONY: build run test tidy fmt fmt-check vet lint arch-lint ci install-tools clean docker-build docker-run help


build: ## Собрать бинарник в bin/
	go build -o bin/$(BINARY) $(CMD)

run: ## Запустить утилиту локально
	go run $(CMD)

test: ## Прогнать тесты
	go test -race ./...

tidy: ## Привести в порядок зависимости
	go mod tidy

fmt: ## Отформатировать код
	gofmt -w .

fmt-check: ## Проверить форматирование
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Файлы не отформатированы (запусти make fmt):"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet: ## Статический анализ go vet
	go vet ./...

lint: ## Запустить golangci-lint
	$(GOLANGCI_LINT) run --timeout=5m

arch-lint: ## Проверить архитектурные границы
	$(GO_ARCH_LINT) check

ci: fmt-check lint test arch-lint ## Полный набор проверок CI

install-tools: ## Установить линтеры (golangci-lint, go-arch-lint)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOBIN)
	go install github.com/fe3dback/go-arch-lint@latest

clean: ## Удалить артефакты сборки
	rm -rf bin

docker-build: ## Собрать Docker-образ
	docker build -t $(IMAGE) .

docker-run: ## Запустить утилиту в контейнере (интерактивно)
	docker run --rm -it --env-file .env -v $(PWD)/downloads:/app/downloads $(IMAGE)

help: ## Показать список целей
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

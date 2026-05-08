NAME=go-echo-skeleton
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

default: all

all:
	@go generate ./internal/enum
	golangci-lint run
lint:
	golangci-lint run
generate:
	@go generate ./internal/enum
docs:
	@swag init

# ── code generation ──────────────────────────────────

# 从数据库反向生成 models (DSN 从 config/config.yaml 读取)
# 用法: make gen-models                      ← 全部表
#       make gen-models TABLES="order user"  ← 只生成指定表
#       make gen-models-file                 ← 直接使用 generate/xorm/reverse.yml
gen-models:
	go run ./cmd/gen models $(if $(TABLES),--tables $(TABLES))
gen-models-file:
	reverse -f generate/xorm/reverse.yml

# 为指定模型生成 repo + service 脚手架
# 用法: make gen-crud Models="User Order Ticket"
gen-crud:
	go run ./cmd/gen crud  $(if $(Models),--models $(Models))

# ── build ────────────────────────────────────────────

win64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o ./bin/$(NAME).exe ./cmd/server
	@echo "build success"
	
linux-amd64:
	@rm -f $(NAME)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o ./bin/$(NAME) ./cmd/server
	@chmod +w bin/$(NAME)
	@upx -9 bin/$(NAME)
	@echo "build success"
	
.PHONY: all lint generate docs gen-models gen-models-file gen-crud win64 linux-amd64 

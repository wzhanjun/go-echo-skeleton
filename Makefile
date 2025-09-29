NAME=go-echo-skeleton
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'

default: services all

all:
	@go generate ./enum
	golangci-lint run 
lint:
	golangci-lint run 
generate:
	@go generate ./enum
docs:
	@swag init	

win64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o ./bin/$(NAME).exe ./cmd/server
	@echo "build success"
	
linux-amd64:
	@rm -f $(NAME)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o ./bin/$(NAME) ./cmd/server
	@chmod +w bin/$(NAME)
	@upx -9 bin/$(NAME)
	@echo "build success"
	
.PHONY: all lint generate docs 

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
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(NAME).exe
	@echo "build success"
	
linux-amd64:
	@rm -f $(NAME)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(NAME)
	@chmod +w $(NAME)
	@upx -9 $(NAME)
	@echo "build success"
	
.PHONY: all lint generate docs 

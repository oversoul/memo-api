run:
	@go run cmd/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -o go-app cmd/main.go

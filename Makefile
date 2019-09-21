vet:
	go vet ./...
fmt:
	go fmt ./...

build: vet fmt
	 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app main.go
	 docker build -t magicsong/recommend bin/
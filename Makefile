build:
	go build -o out/notifier ./...

run:
	go run ./... -config example.yaml

test:
	go test ./...

lint:
	golangci-lint run

racecondition:
	go test -race -short ./...

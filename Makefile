test:
	go test ./...

test-e2e:
	go test --tags=e2e ./...

test-integration:
	go test --tags=integration ./...

install:
	go get

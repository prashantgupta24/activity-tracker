example-start:
	go run example/example.go

test: vet
	go test -v -failfast -race ./...

vet:
	go vet $(shell glide nv)

lint:
	go list ./... | xargs -L1 golint -set_exit_status
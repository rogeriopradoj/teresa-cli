all:
	@godep go build -o teresa

test:
	@godep go test -v ./cmd

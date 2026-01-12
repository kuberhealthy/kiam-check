IMAGE := "kuberhealthy/kiam-check"
TAG := "latest"

# Build the KIAM check container locally.
build:
	podman build -f Containerfile -t {{IMAGE}}:{{TAG}} .

# Run the unit tests for the KIAM check.
test:
	go test ./...

# Build the KIAM check binary locally.
binary:
	go build -o bin/kiam-check ./cmd/kiam-check

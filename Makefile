all: vendor fmt vet test build

vendor: glide.lock
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		-w /go/src/github.com/trusch/authy \
		golang:1.9 bash -c "\
			curl https://glide.sh/get | sh;\
			glide install;\
		"

fmt: $(shell find ./cmd ./http)
	test $(shell docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		golang:1.9 go fmt github.com/trusch/authy/...| wc -l) -eq 0

vet: $(shell find ./cmd ./http)
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		golang:1.9 go vet github.com/trusch/authy/...

test: $(shell find ./cmd ./http)
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		golang:1.9 go test -cover github.com/trusch/authy/...

build: $(shell find ./cmd ./http)
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		-w /go/src/github.com/trusch/authy \
		-e CGO_ENABLED=0 \
		golang:1.9 go build -v --ldflags '-extldflags "-static"'

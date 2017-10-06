all: vendor fmt vet test authy image

vendor: glide.lock
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		-w /go/src/github.com/trusch/authy \
		instrumentisto/glide:0.13 --home /tmp install

glide.lock: glide.yaml
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		-w /go/src/github.com/trusch/authy \
		instrumentisto/glide:0.13 --home /tmp up

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

authy: $(shell find ./cmd ./http)
	docker run \
		-u $(shell stat -c "%u:%g" .) \
		-v $(shell pwd):/go/src/github.com/trusch/authy \
		-w /go/src/github.com/trusch/authy \
		-e CGO_ENABLED=0 \
		golang:1.9 go build -v --ldflags '-extldflags "-static"'

image: authy
	cp authy docker
	cd docker && docker build -t trusch/authy .

clean:
	rm -rf authy vendor glide.lock

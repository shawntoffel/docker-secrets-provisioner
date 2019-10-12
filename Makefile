TAG=$(shell git describe --tags --always)
VERSION=$(TAG:v%=%)
NAME=dsp
REPO=shawntoffel/$(NAME)
GO=GO111MODULE=on go
BUILD=GOARCH=amd64 $(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" 

.PHONY: all deps test build build-linux docker-build docker-save docker-deploy clean 

all: deps test build 
deps:
	$(GO) mod download

test:
	$(GO) vet ./...
	$(GO) test -v -race ./...

build:
	$(BUILD) -o bin/$(NAME)-$(VERSION) ./cmd/...

build-linux:
	CGO_ENABLED=0 GOOS=linux $(BUILD) -a -installsuffix cgo -o bin/$(NAME) ./cmd/...

docker-build:
	docker build -t $(REPO):$(VERSION) .

docker-save:
	mkdir -p bin && docker save -o bin/image.tar $(REPO):$(VERSION)

docker-deploy:
	docker push $(REPO):$(VERSION)

clean:
	@find bin -type f ! -name '*.toml' -delete -print
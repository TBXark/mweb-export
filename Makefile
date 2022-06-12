NAME=mweb-export
BINDIR=bin
GOBUILD=go build
VERSION=$(shell git describe --tags || echo "unknown version")

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \

all: darwin-amd64 darwin-arm64

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))

$(gz_releases): %.gz : %
	chmod +x $(BINDIR)/$(NAME)-$(basename $@)
	gzip -f -S -$(VERSION).gz $(BINDIR)/$(NAME)-$(basename $@)

all-arch: $(PLATFORM_LIST)

releases: $(gz_releases)

lint:
	golangci-lint run --disable-all -E govet -E gofumpt -E megacheck ./...

clean:
	rm $(BINDIR)/*

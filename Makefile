ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=SET
    NUL=NUL
else
    SET=export
    NUL=/dev/null
endif

NAME=$(notdir $(CURDIR))
VERSION=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
EXE:=$(shell go env GOEXE)

all:
	go fmt ./...
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

_dist:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXE)

dist:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist

release:
	$(go) run github.com/hymkor/latest-notes@master | gh release create -d --notes-file - -t $(version) $(version) $(wildcard $(name)-$(version)-*.zip)

manifest:
	$(GO) run github.com/hymkor/make-scoop-manifest@master -all *-windows-*.zip > $(NAME).json

.PHONY: dist _dist manifest all release manifest

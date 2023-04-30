ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=SET
    NUL=NUL
else
    SET=export
    NUL=/dev/null
endif

NAME=$(lastword $(subst /, ,$(abspath .)))
VERSION=$(shell git.exe describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"

all:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

_package:
	$(MAKE) all
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXT)

package:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _package EXT=
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _package EXT=
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _package EXT=.exe
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _package EXT=.exe

release:
	gh release create -d --notes "" -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

manifest:
	make-scoop-manifest *-windows-*.zip > $(NAME).json

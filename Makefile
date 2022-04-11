export GO111MODULE=on
VERSION="$(shell git describe --tags --candidates=1 --dirty)+envato"
FLAGS=-X main.Version=$(VERSION) -s -w
SOURCES=$(wildcard *.go iamy/*.go)

# To create a new release:
#  $ git tag vx.x.x
#  $ git push --tags
#  $ make clean
#  $ make release     # this will create 3 binaries in ./bin
#
#  Next, go to https://github.com/99designs/iamy/releases/new
#  - select the tag version you just created
#  - Attach the binaries from ./bin/*

release: bin/iamy-linux-arm64 bin/iamy-linux-amd64 bin/iamy-darwin-amd64 bin/iamy-windows-386.exe bin/iamy-darwin-arm64 bin/iamy-freebsd-amd64

bin/iamy-darwin-arm64: $(SOURCES)
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -o $@ -ldflags="$(FLAGS)" .

bin/iamy-linux-arm64: $(SOURCES)
	@mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -o $@ -ldflags="$(FLAGS)" .

bin/iamy-linux-amd64: $(SOURCES)
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

bin/iamy-darwin-amd64: $(SOURCES)
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

bin/iamy-windows-386.exe: $(SOURCES)
	@mkdir -p bin
	GOOS=windows GOARCH=386 go build -o $@ -ldflags="$(FLAGS)" .

bin/iamy-freebsd-amd64: $(SOURCES)
	@mkdir -p bin
	GOOS=freebsd GOARCH=amd64 go build -o $@ -ldflags="$(FLAGS)" .

clean:
	rm -f bin/*

.PHONY: clean release

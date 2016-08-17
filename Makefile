GO_EXECUTABLE := go
VERSION := $(shell git describe --abbrev=10 --dirty --always --tags)
DIST_DIRS := find * -type d -exec

all: build install

build:
	${GO_EXECUTABLE} build -o gotree -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree

install:
	${GO_EXECUTABLE} install -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree

test:
	${GO_EXECUTABLE} test github.com/fredericlemoine/gotree/tests/

deploy:
	mkdir -p deploy
	env GOOS=windows GOARCH=amd64 ${GO_EXECUTABLE} build -o deploy/gotree_amd64.exe -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree
	env GOOS=windows GOARCH=386 ${GO_EXECUTABLE} build -o deploy/gotree_386.exe -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree
	env GOOS=darwin GOARCH=amd64 ${GO_EXECUTABLE} build -o deploy/gotree_amd64_darwin -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree
	env GOOS=darwin GOARCH=386 ${GO_EXECUTABLE} build -o deploy/gotree_386_darwin -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree
	env GOOS=linux GOARCH=amd64 ${GO_EXECUTABLE} build -o deploy/gotree_amd64_linux -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree
	env GOOS=linux GOARCH=386 ${GO_EXECUTABLE} build -o deploy/gotree_386_linux -ldflags "-X github.com/fredericlemoine/gotree/cmd.Version=${VERSION}" github.com/fredericlemoine/gotree

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

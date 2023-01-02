GO_EXECUTABLE := go
VERSION := $(shell git describe --abbrev=10 --dirty --always --tags)
VERSION_PACKAGE := github.com/evolbioinfo/gotree/cmd.Version
NAME := gotree
PACKAGE:=github.com/evolbioinfo/gotree
CGO_ENABLED:=0

all: dep build test

dep:
	${GO_EXECUTABLE} get .

build:
	CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} build -o ${NAME} -ldflags "-X ${VERSION_PACKAGE}=${VERSION}" ${PACKAGE}

install: testcommands
	rm -f ${GOPATH}/bin/${NAME}
	CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} install -ldflags "-X ${VERSION_PACKAGE}=${VERSION}" ${PACKAGE}

test: dep
	CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} test ${PACKAGE}/...

testcommands:
	bash ./test.sh

.PHONY: deploy deploydir deploywinamd deploywin386 deploylinuxamd deploylinux386 deploydarwinamd deploydarwin386

deploy: deploywinamd deploywin386 deploylinuxamd deploylinux386 deploydarwinamd deploydarwinarm
	tar -czvf deploy/${VERSION}.tar.gz --directory="deploy" ${VERSION}

deploydir:
	mkdir -p deploy/${VERSION}

deploywinamd: deploydir dep
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} build -o deploy/${VERSION}/${NAME}_amd64.exe -ldflags "-X ${VERSION_PACKAGE}=${VERSION}" ${PACKAGE}

deploylinuxamd: deploydir dep
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} build -o deploy/${VERSION}/${NAME}_amd64_linux -ldflags "-X ${VERSION_PACKAGE}=${VERSION}" ${PACKAGE}

deploydarwinamd: deploydir dep
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} build -o deploy/${VERSION}/${NAME}_amd64_darwin -ldflags "-X ${VERSION_PACKAGE}=${VERSION}" ${PACKAGE}

deploydarwinarm: deploydir dep
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=${CGO_ENABLED} ${GO_EXECUTABLE} build -o deploy/${VERSION}/${NAME}_arm64_darwin -ldflags "-X ${VERSION_PACKAGE}=${VERSION}" ${PACKAGE}

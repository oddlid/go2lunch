BINARY := lhlunch.bin
VERSION := 2019-01-22
UNAME := $(shell uname -s)
SOURCES := $(wildcard *.go)
COMMIT_ID := $(shell git describe --tags --always)
BUILD_TIME := $(shell go run tool/rfc3339date.go)
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.BUILD_DATE=${BUILD_TIME} -X main.COMMIT_ID=${COMMIT_ID} -s -w ${DFLAG}"

ifeq ($(UNAME), Linux)
	DFLAG := -d
endif

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	env CGO_ENABLED=0 go build ${LDFLAGS} -o $@ ${SOURCES}

.PHONY: install
install:
	env CGO_ENABLED=0 go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm -f ${BINARY}; fi

#.PHONY: dbg
#dbg:
#	echo env CGO_ENABLED=0 go build ${LDFLAGS} -o $@ ${SOURCES}

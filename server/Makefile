BINARY := go2lunchsrv.bin
VERSION := 2025-02-23
UNAME := $(shell uname -s)
SOURCES := $(wildcard *.go)
COMMIT_ID := $(shell git describe --tags --always)
BUILD_TIME := $(shell go run tool/rfc3339date.go)
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.BuildDate=${BUILD_TIME} -X main.CommitID=${COMMIT_ID} -s -w ${DFLAG}"

ifeq ($(UNAME), Linux)
	DFLAG := -d
endif

.DEFAULT_GOAL: $(BINARY)

# Since we have build constraints, we should pass '.' (package) to build, not a list of go files
$(BINARY): $(SOURCES)
	env CGO_ENABLED=0 go build ${LDFLAGS} -o $@ .

.PHONY: install
install:
	env CGO_ENABLED=0 go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm -f ${BINARY}; fi

#.PHONY: dbg
#dbg:
#	echo env CGO_ENABLED=0 go build ${LDFLAGS} -o $@ ${SOURCES}

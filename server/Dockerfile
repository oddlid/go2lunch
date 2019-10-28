FROM golang:stretch as builder

ARG BIN=go2lunchsrv.bin

# Note that $GOPATH is set to /go from the parent image
WORKDIR ${GOPATH}/src/github.com/oddlid/
RUN git clone https://github.com/oddlid/go2lunch.git
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch
RUN git checkout dev
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch/server
RUN go get -d -v ./...
# This needs to be done in addition, for the binary rice. Will not get pulled in from the step above
RUN go get -u -v github.com/GeertJohan/go.rice/rice
RUN make
RUN ${GOPATH}/bin/rice append --exec ${BIN}


FROM alpine:latest
LABEL maintainer="Odd E. Ebbesen <oddebb@gmail.com>"
RUN apk add --no-cache --update tini ca-certificates \
		&& \
		rm -rf /var/cache/apk/*
RUN adduser -D -u 1000 lunchsrv
COPY --from=builder /go/src/github.com/oddlid/go2lunch/server/${BIN} /usr/local/bin/
RUN chown lunchsrv /usr/local/bin/${BIN} && chmod 555 /usr/local/bin/${BIN}
EXPOSE 20666 20667
USER lunchsrv
ENTRYPOINT ["tini", "-g", "--", ${BIN}]
CMD ["-h"]

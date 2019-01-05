FROM golang:stretch as builder

# Note that $GOPATH is set to /go from the parent image
WORKDIR ${GOPATH}/src/github.com/oddlid/
RUN git clone https://github.com/oddlid/go2lunch.git
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch
RUN git checkout dev
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch/lhlunch
RUN go get -d -v ./...
# This needs to be done in addition, for the binary rice. Will not get pulled in from the step above
RUN go get -u -v github.com/GeertJohan/go.rice/rice
RUN make
RUN ${GOPATH}/bin/rice append --exec lhlunch.bin


FROM alpine:latest
MAINTAINER Odd E. Ebbesen <oddebb@gmail.com>
RUN apk add --no-cache --update tini ca-certificates \
		&& \
		rm -rf /var/cache/apk/*
RUN adduser -D -u 1000 lunchsrv
COPY --from=builder /go/src/github.com/oddlid/go2lunch/lhlunch/lhlunch.bin /usr/local/bin/lhlunch
RUN chown lunchsrv /usr/local/bin/lhlunch && chmod 555 /usr/local/bin/lhlunch
EXPOSE 3000
ENV PATH="/usr/local/bin:${PATH}"
USER lunchsrv
ENTRYPOINT ["tini", "-g", "--", "lhlunch"]
CMD ["serve", "-l", ":3000", "--writepid", "/tmp/lhlunch.pid"]

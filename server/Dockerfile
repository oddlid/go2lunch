ARG BUILD_BIN=go2lunchsrv.bin

FROM golang:stretch as builder
# Specifying the same ARG without value copies it into this FROM-scope, from the outside scoped ARG before FROM
ARG BUILD_BIN
# Setting TZ makes "make" build the correct time (at least for local builds) into the executable
ARG ARG_TZ=Europe/Stockholm
ENV TZ=${ARG_TZ}
# Note that $GOPATH is set to /go from the parent image
WORKDIR ${GOPATH}/src/github.com/oddlid/
RUN git clone https://github.com/oddlid/go2lunch.git
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch
RUN git checkout dev
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch/server
RUN go get -d -v ./...
# This needs to be done in addition, for the binary rice. Will not get pulled in from the step above
RUN go get -u -v github.com/GeertJohan/go.rice/rice
RUN make BINARY=${BUILD_BIN}
RUN ${GOPATH}/bin/rice append --exec ${BUILD_BIN}


# We don't do anything about TZ in the image for running the app, as one would rather specify
# that at runtime, and not have it baked into the image

FROM alpine:latest
LABEL maintainer="Odd E. Ebbesen <oddebb@gmail.com>"
ARG BUILD_BIN
RUN apk add --no-cache --update tini ca-certificates \
		&& \
		rm -rf /var/cache/apk/*
RUN adduser -D -u 1000 lunchsrv
COPY --from=builder /go/src/github.com/oddlid/go2lunch/server/${BUILD_BIN} /usr/local/bin/
RUN chown lunchsrv /usr/local/bin/${BUILD_BIN} && chmod 555 /usr/local/bin/${BUILD_BIN}
EXPOSE 20666 20667
USER lunchsrv
ENTRYPOINT ["tini", "-g", "--"]
ENV LUNCHBIN=${BUILD_BIN}
CMD ${BUILD_BIN} -h

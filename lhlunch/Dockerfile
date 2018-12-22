FROM oddlid/arch-cli:latest as builder 
MAINTAINER Odd E. Ebbesen <oddebb@gmail.com>
ENV GOPATH=/go
RUN go get -d -u -v github.com/PuerkitoBio/goquery
RUN go get -d -u -v github.com/robfig/cron
RUN go get -d -u -v github.com/urfave/cli
RUN go get -d -u -v github.com/Sirupsen/logrus
RUN go get -d -u -v github.com/gorilla/mux
RUN go get -d -u -v github.com/GeertJohan/go.rice
RUN go get -u -v github.com/GeertJohan/go.rice/rice
WORKDIR ${GOPATH}/src/github.com/oddlid/
RUN git clone https://github.com/oddlid/go2lunch.git
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch
RUN git checkout dev
WORKDIR ${GOPATH}/src/github.com/oddlid/go2lunch/lhlunch
RUN make
RUN ${GOPATH}/bin/rice append --exec lhlunch.bin

FROM alpine:latest
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

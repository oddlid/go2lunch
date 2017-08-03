FROM alpine:latest
MAINTAINER Odd E. Ebbesen <oddebb@gmail.com>

RUN apk add --update tini ca-certificates
RUN adduser -D -u 1000 lunchsrv

COPY lhlunch.bin /usr/local/bin/lhlunch
RUN chown lunchsrv /usr/local/bin/lhlunch && chmod 555 /usr/local/bin/lhlunch
EXPOSE 3000
USER lunchsrv
ENTRYPOINT ["tini", "-g", "--", "/usr/local/bin/lhlunch"]
CMD ["serve", "-l", ":3000", "--writepid", "/tmp/lhlunch.pid"]

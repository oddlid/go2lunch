FROM alpine:latest
MAINTAINER Odd E. Ebbesen <oddebb@gmail.com>

RUN apk add --update tini ca-certificates

COPY lhlunch.bin /usr/local/bin/lhlunch
RUN chmod 555 /usr/local/bin/lhlunch
EXPOSE 3000
ENTRYPOINT ["tini", "-g", "--", "/usr/local/bin/lhlunch"]
CMD ["serve", "-l", ":3000", "--writepid", "/tmp/lhlunch.pid"]

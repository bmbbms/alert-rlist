FROM golang:alpine3.18
MAINTAINER jlpay.com

ADD ./ /app/

WORKDIR /app

RUN go env -w GO111MODULE=on; \
    go env -w GOPROXY=https://goproxy.cn,direct
RUN set -eux; \
        sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories;\
        apk update;\
        apk upgrade; \
        apk add --no-cache make;\
         ls -la \
        ; make \
        ;./alert-rlist --version



FROM alpine:3.18

COPY --from=0 /app/alert-rlist /usr/bin/

RUN     set -eux ;\
        sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories;\
        apk update;\
        apk upgrade; \
        apk add --no-cache tini \
        \
        ; chmod +x /usr/bin/alert-rlist \
        ; /usr/bin/alert-rlist --version


ENTRYPOINT ["/sbin/tini", "--"]
CMD /usr/bin/alert-rlist

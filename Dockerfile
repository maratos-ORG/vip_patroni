FROM golang:1.17

ARG app_version
ENV APP /go/src/vip_patroni

ENV \
  CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

RUN mkdir -p $APP
RUN mkdir -p /root/build

WORKDIR $APP
COPY . $APP/src
WORKDIR $APP/src

RUN go mod download
RUN go mod verify

RUN make test
#RUN make build
RUN make build RELEASE=${app_version}

RUN mv $APP/src/bin/vip_patroni* /root/build

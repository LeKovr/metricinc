
ARG golang_version

FROM golang:$golang_version

MAINTAINER Alexey Kovrizhkin <lekovr+docker@gmail.com>

WORKDIR /go/src/github.com/LeKovr/metricinc
COPY cmd cmd
COPY lib lib
COPY counter counter
COPY Makefile .
COPY glide.* ./

RUN go get -u github.com/golang/lint/golint
RUN make vendor
RUN make build-standalone

FROM scratch

VOLUME /data

WORKDIR /
COPY --from=0 /go/src/github.com/LeKovr/metricinc/client .
COPY --from=0 /go/src/github.com/LeKovr/metricinc/server .

EXPOSE 50051
CMD ["/server"]

FROM golang:1.22

WORKDIR ${GOPATH}/merchStore/
COPY . ${GOPATH}/merchStore/

RUN go build -o /build ./cmd \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]
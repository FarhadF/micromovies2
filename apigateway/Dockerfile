FROM golang:1.10.3-alpine3.7 as builder

RUN apk update; apk add git curl && go get -u github.com/golang/dep/cmd/dep
COPY ["./Gopkg.lock", "./Gopkg.toml", "/go/src/micromovies2/apigateway/"]
RUN cd /go/src/micromovies2/apigateway && dep ensure --vendor-only
COPY ./ /go/src/micromovies2/apigateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo ./cmd/server.go

FROM scratch
COPY --from=builder /go/src/micromovies2/apigateway/server /
ENTRYPOINT ["/server"]
FROM golang:1.10.3-alpine3.7 as builder

RUN apk update; apk add git curl && go get -u github.com/golang/dep/cmd/dep
COPY ./ /go/src/github.com/farhadf/micromovies2/apigateway
WORKDIR /go/src/github.com/farhadf/micromovies2/apigateway
RUN cd /go/src/github.com/farhadf/micromovies2/apigateway && dep ensure

RUN ls && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo /go/src/github.com/farhadf/micromovies2/apigateway/cmd/server.go

FROM scratch
COPY --from=builder /go/src/github.com/farhadf/micromovies2/apigateway/server /
ENTRYPOINT ["/server"]
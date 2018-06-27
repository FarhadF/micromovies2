FROM golang:1.10.3-alpine3.7 as builder

RUN apk update; apk add git curl && go get -u golang.org/x/vgo && mkdir /go/src/github.com/farhadf/ -p && \
    cd /go/src/github.com/farhadf/ && git clone https://github.com/farhadf/micromovies2

#COPY ./ /go/src/github.com/farhadf/micromovies2/apigateway
WORKDIR /go/src/github.com/farhadf/micromovies2/apigateway
RUN cd /go/src/github.com/farhadf/micromovies2/apigateway  && \
CGO_ENABLED=0 GOOS=linux vgo build -a -installsuffix nocgo /go/src/github.com/farhadf/micromovies2/apigateway/cmd/server.go
FROM scratch
COPY --from=builder ["/go/src/github.com/farhadf/micromovies2/apigateway/server", "/go/src/github.com/farhadf/micromovies2/apigateway/cmd/model.conf", "/go/src/github.com/farhadf/micromovies2/apigateway/cmd/policy.csv", "/"]
ENTRYPOINT ["/server"]
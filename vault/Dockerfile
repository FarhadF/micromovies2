FROM golang:1.10.3-alpine3.7 as builder

RUN apk update; apk add git curl && go get -u golang.org/x/vgo && mkdir /go/src/github.com/farhadf/ -p && \
    cd /go/src/github.com/farhadf/ && git clone https://github.com/farhadf/micromovies2
#COPY ./ /go/src/github.com/farhadf/micromovies2/movies
WORKDIR /go/src/github.com/farhadf/micromovies2/movies
RUN cd /go/src/github.com/farhadf/micromovies2/movies  && \
CGO_ENABLED=0 GOOS=linux vgo build -a -installsuffix nocgo /go/src/github.com/farhadf/micromovies2/movies/cmd/server.go
FROM scratch
COPY --from=builder ["/go/src/github.com/farhadf/micromovies2/movies/server", "/"]
ENTRYPOINT ["/server"]
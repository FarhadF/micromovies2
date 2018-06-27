FROM golang:1.10.3-alpine3.7 as builder

RUN apk update; apk add git curl && go get -u golang.org/x/vgo
COPY ./ /go/src/github.com/farhadf/micromovies2/jwtauth
WORKDIR /go/src/github.com/farhadf/micromovies2/jwtauth
RUN cd /go/src/github.com/farhadf/micromovies2/jwtauth  && \
CGO_ENABLED=0 GOOS=linux vgo build -a -installsuffix nocgo /go/src/github.com/farhadf/micromovies2/jwtauth/cmd/server.go
FROM scratch
COPY --from=builder ["/go/src/github.com/farhadf/micromovies2/jwtauth/server", "/"]
CMD ["/server"]
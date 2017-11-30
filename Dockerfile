FROM golang:1.9 AS gobuild

RUN go get -v github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/dweomer/sweomer

WORKDIR /go/src/github.com/dweomer/sweomer

ARG VERSION
ENV CGO_ENABLED=0

RUN dep ensure -v
RUN export VERSION=${VERSION:-$(git describe --tags --always --dirty | sed -e 's/^v//g' -e "s/dirty/dev-$(git rev-parse --short HEAD)/g")} \
 && go install -v -ldflags="-X main.Version=${VERSION}" ./...
RUN /go/bin/sweomer --version

FROM scratch

COPY --from=gobuild /go/bin/sweomer /

ENTRYPOINT ["/sweomer"]
CMD ["help"]

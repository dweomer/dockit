FROM golang:1.9 AS gobuild

RUN go get -v github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/dweomer/sweomer

WORKDIR /go/src/github.com/dweomer/sweomer

ENV CGO_ENABLED=0

RUN dep ensure -v
RUN go install -v ./...

FROM scratch

COPY --from=gobuild /go/bin/sweomer /

ENTRYPOINT ["/sweomer"]
CMD ["help"]

FROM ethereum/solc:0.8.17-alpine AS base

COPY --from=golang:alpine3.16 /usr/local/go/ /usr/local/go/
ENV PATH /usr/local/go/bin:$PATH

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN apk add --no-cache make gcc musl-dev tini

RUN go install github.com/ethereum/go-ethereum/cmd/abigen@v1.10.25

WORKDIR /go/src/bazaar

COPY go.mod go.mod
RUN go mod download

COPY . .

RUN make sol
RUN make build

ENTRYPOINT ["/sbin/tini", "-g", "--"]
CMD ["/go/src/bazaar/scripts/run.sh"]

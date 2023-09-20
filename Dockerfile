FROM quay.io/goswagger/swagger:latest as swagger

FROM golang:1.18-alpine as builder

RUN apk add --no-cache ca-certificates git mercurial

ENV GOPRIVATE=github.com

WORKDIR /dist

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY pokeNames.txt adj.txt capabilities.txt ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-s -w" \
    -o pokemon-factory .

COPY --from=swagger /usr/bin/swagger /usr/bin/swagger
RUN swagger generate spec -m -o ./swagger.json

FROM alpine:3.14

RUN apk update update-ca-certificates
COPY --from=builder /dist/pokemon-factory /go/bin/
COPY --from=builder /dist/swagger.json /swagger.json
COPY --from=builder /dist/pokeNames.txt /dist/adj.txt /dist/capabilities.txt /

ENTRYPOINT [ "/go/bin/pokemon-factory" ]

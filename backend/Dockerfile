FROM golang:1.19 AS build
WORKDIR /go/src
COPY model ./go
COPY main.go .
COPY go.sum .
COPY go.mod .

ENV CGO_ENABLED=0

RUN go build -o openapi .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /model/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]

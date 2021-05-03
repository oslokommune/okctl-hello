FROM golang:1.16 AS build
ENV CGO_ENABLED=0
WORKDIR /go/src

COPY go.* ./
RUN go get -d -v ./...

COPY ./pkg ./pkg
COPY main.go .
COPY ./public ./public

RUN go build -a -installsuffix cgo -o openapi .

FROM scratch AS runtime
ENTRYPOINT ["./openapi"]
EXPOSE 3000/tcp

COPY --from=build /go/src/openapi ./

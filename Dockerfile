# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY webapp/go.mod .
COPY webapp/go.sum .

RUN go mod download

COPY webapp/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# final stage
FROM scratch
COPY --from=builder /app/webapp /app/
EXPOSE 8080
ENTRYPOINT ["/app/webapp"]
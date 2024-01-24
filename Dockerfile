FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/tea

FROM scratch

WORKDIR /app

COPY --from=builder /app/tea ./

CMD ["./tea"]
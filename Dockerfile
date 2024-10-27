FROM golang:1.23-alpine AS builder
LABEL authors="varankin_vasilii"

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
COPY ./internal/server/storage/db/migrations ./internal/storage/db/migrations/
RUN go build -o ./bin/keeper cmd/server/main.go

FROM alpine as runner
COPY --from=builder /usr/local/src/bin/keeper /

CMD ["./keeper"]
EXPOSE 3200
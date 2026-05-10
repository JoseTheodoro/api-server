FROM golang:1.26 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download 

COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o api ./cmd/api

FROM alpine:3.23.3
WORKDIR /app/
COPY --from=builder /build/api .
EXPOSE 9090
CMD [ "./api" ]







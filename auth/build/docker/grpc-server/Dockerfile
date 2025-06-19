# Build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /opt/auth
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build -o bin/app ./cmd/grpc_server/main.go

# Final stage
FROM alpine:latest

WORKDIR /opt/auth

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Europe/Moscow
RUN cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone
    
COPY --from=builder /opt/auth/bin/app .

RUN chmod +x ./app

CMD [ "./app" ]

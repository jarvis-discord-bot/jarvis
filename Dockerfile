FROM golang:1.21-alpine AS builder

WORKDIR /

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o jarvis /cmd/jarvis/main.go

RUN apk --update add ca-certificates

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder /jarvis .
COPY --from=builder /db/migrations /db/migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV JARVIS_PORT=${JARVIS_PORT:-8080}
ENV JARVIS_ADDRESS="0.0.0.0"
ENV JARVIS_SQL_MIGRATION_PATH="/db/migrations"

# Export necessary port.
EXPOSE ${JARVIS_PORT}

# Command to run when starting the container.
ENTRYPOINT ["./jarvis"]
FROM klakegg/hugo:ext-alpine as hugo-builder

COPY . /src
WORKDIR /src

RUN make build

# Stage 2: Build the Go binary
FROM golang:1.23-alpine as go-builder

# Copy the Go source files
COPY . /src
WORKDIR /src

# Build the Go binary
RUN go build -o /app/main .

FROM alpine:latest

COPY --from=hugo-builder /src/public /app/public
COPY --from=hugo-builder /src/templates /app/templates
COPY --from=go-builder /app/main /app/main

WORKDIR /app

# Expose the port the app runs on
EXPOSE 8080

# Run the Go binary
CMD ["/app/main"]
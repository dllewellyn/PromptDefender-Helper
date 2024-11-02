FROM klakegg/hugo:latest as hugo-builder

COPY . /src
WORKDIR /src

RUN cd ui && hugo && mv public ../public

FROM golang:1.23-alpine as go-builder

COPY . /src
WORKDIR /src

# Build the Go binary
RUN go build -o /app/main .

FROM alpine:latest

COPY --from=hugo-builder /src/public /app/public
COPY --from=hugo-builder /src/templates /app/templates
COPY --from=go-builder /app/main /app/main
COPY --from=go-builder /src/prompts /app/prompts

WORKDIR /app

# Expose the port the app runs on
EXPOSE 8080

# Run the Go binary
CMD ["sh", "-c", "cd /app && ./main"]
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
ENV PORT 8080

EXPOSE 8080

ENV GOOGLE_APPLICATION_CREDENTIALS="/app/service-account.json"

COPY entrypoint.sh .

RUN chmod +x /app/entrypoint.sh

# Run the Go binary
ENTRYPOINT ["/app/entrypoint.sh"]
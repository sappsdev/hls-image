# Stage 1: Build the go static binary
FROM golang:1.17.5-alpine3.15 AS server-builder
RUN apk update && apk upgrade && \
  apk --update add git
WORKDIR /builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' -a \
  -o server

# Stage 2: Final
FROM alpine
WORKDIR /app
COPY --from=server-builder /builder/server .
COPY --from=server-builder /builder/fonts ./fonts/.
RUN apk add ffmpeg

EXPOSE 4000
ENTRYPOINT ["./server"]
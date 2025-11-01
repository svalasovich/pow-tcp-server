# Build stage
FROM golang:1.25 AS builder

ENV GOOS=linux

WORKDIR /src
COPY . /src

RUN make build

FROM gcr.io/distroless/static:latest

COPY --from=builder /src/dist/server /usr/local/bin

ENTRYPOINT ["server"]

EXPOSE 9000

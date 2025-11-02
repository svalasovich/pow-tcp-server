# Build stage
FROM golang:1.25 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /src
COPY . /src

RUN make build

FROM gcr.io/distroless/static:latest

COPY --from=builder /src/dist /usr/local/bin

ENTRYPOINT ["client"]

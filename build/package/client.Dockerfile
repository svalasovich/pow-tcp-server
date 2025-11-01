# Build stage
FROM golang:1.25 AS builder

ENV GOOS=linux

WORKDIR /src
COPY . /src

RUN make build

FROM gcr.io/distroless/static:latest

COPY --from=builder /src/dist/client /usr/local/bin

ENTRYPOINT ["client"]

EXPOSE 9000

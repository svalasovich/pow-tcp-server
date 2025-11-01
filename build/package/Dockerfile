# Build stage
FROM golang:1.23 AS builder

ARG GIT_SHA='N/A'

ENV GOOS=linux
ENV GIT_SHA=$GIT_SHA

WORKDIR /src
COPY . /src

RUN make build

FROM gcr.io/distroless/static:latest

COPY --from=builder /src/dist /usr/local/bin

ENTRYPOINT ["golang-template"]

EXPOSE 9000

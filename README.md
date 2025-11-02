# TCP Server (protected from DDOS attacks with the Proof of Work) & Client

<p align="center">
  <a href="https://github.com/svalasovich/pow-tcp-server/actions/workflows/lint.yml?query=branch%3Amain"><img src="https://github.com/svalasovich/pow-tcp-server/actions/workflows/lint.yml/badge.svg?branch=main"></a>
  <a href="https://github.com/svalasovich/pow-tcp-server/actions/workflows/test.yml?query=branch%3Amain"><img src="https://github.com/svalasovich/pow-tcp-server/actions/workflows/test.yml/badge.svg?branch=main"></a>
  <a href="https://github.com/svalasovich/pow-tcp-server/actions/workflows/build.yml?query=branch%3Amain"><img src="https://github.com/svalasovich/pow-tcp-server/actions/workflows/build.yml/badge.svg?branch=main"></a>
  <br>
  <a href="https://goreportcard.com/report/github.com/svalasovich/pow-tcp-server"><img src="https://goreportcard.com/badge/github.com/svalasovich/pow-tcp-server" /></a>
  <a href="https://github.com/svalasovich/pow-tcp-server/releases/latest"><img src="https://img.shields.io/github/release/svalasovich/pow-tcp-server.svg" /></a>
  <a href="https://go.dev/doc/devel/release#go1.23.0"><img src="https://img.shields.io/badge/golang-%3E%3D1.23.0-blue.svg" /></a>
  <br>
</p>

## Task Description

Design and implement “Word of Wisdom” tcp server:

- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work),
  the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other
  collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## How to build native

```shell
make build
```

## How to run



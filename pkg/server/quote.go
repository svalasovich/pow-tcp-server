package server

import (
	"context"

	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

type (
	QuoteService interface {
		GetRandom() string
	}

	Quote struct {
		service QuoteService
	}
)

func NewQuoteHandler(service QuoteService) *Quote {
	return &Quote{service: service}
}

func (q *Quote) Handle(_ context.Context, conn tcp.Connection) error {
	return conn.Write([]byte(q.service.GetRandom()))
}

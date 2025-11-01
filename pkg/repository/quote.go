package repository

import (
	"bufio"
	"bytes"
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed quotes.txt
var rawQuotes []byte

type Quote struct {
	quotes []string
}

func NewQuote() *Quote {
	reader := bufio.NewReader(bytes.NewReader(rawQuotes))
	var quotes []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			quotes = append(quotes, line)
		}
	}

	return &Quote{
		quotes: quotes[:len(quotes):len(quotes)],
	}
}

func (q *Quote) GetRandom() string {
	return q.quotes[rand.Intn(len(q.quotes))]
}

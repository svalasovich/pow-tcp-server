package service

type (
	QuoteRepo interface {
		GetRandom() string
	}

	Quote struct {
		repo QuoteRepo
	}
)

func NewQuote(repo QuoteRepo) *Quote {
	return &Quote{
		repo: repo,
	}
}

func (q *Quote) GetRandom() string {
	return q.repo.GetRandom()
}

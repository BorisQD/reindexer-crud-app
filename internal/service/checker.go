package service

import "context"

type client interface {
	IsConnected(ctx context.Context) bool
}
type Checker struct {
	db client
}

func NewChecker(db client) *Checker {
	return &Checker{db}
}

func (c Checker) HealthCheck(ctx context.Context) bool {
	return c.db.IsConnected(ctx)
}

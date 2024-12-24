package domain

import (
	"context"
	"time"
)

type Contract struct {
	Id, Description string
	Amount          float64
	Periods         int
	Date            time.Time
	Payments        []Payment
}

type Payment struct {
	Id     string
	Amount float64
	Date   time.Time
}

type ContractRepository interface {
	List(ctx context.Context) ([]Contract, error)
}

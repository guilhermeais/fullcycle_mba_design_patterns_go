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

type InvoiceType string

const (
	InvoiceTypeCash    InvoiceType = "cash"
	InvoiceTypeAccrual InvoiceType = "accrual"
)

func (c Contract) GenerateInvoices(month, year int, invoiceType InvoiceType) ([]Invoice, error) {
	strategy, err := MakeInvoiceGenerationStrategy(invoiceType)

	if err != nil {
		return nil, err
	}

	return strategy.Generate(c, month, year), nil
}

type Payment struct {
	Id     string
	Amount float64
	Date   time.Time
}

type Invoice struct {
	Date   time.Time
	Amount float64
}
type ContractRepository interface {
	List(ctx context.Context) ([]Contract, error)
}

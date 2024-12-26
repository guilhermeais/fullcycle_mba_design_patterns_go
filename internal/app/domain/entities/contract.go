package domain

import (
	"context"
	"fmt"
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
	strategy, err := makeInvoiceGenerationStrategy(invoiceType)

	if err != nil {
		return nil, err
	}

	return strategy.Generate(c, month, year), nil
}

func makeInvoiceGenerationStrategy(invoiceType InvoiceType) (InvoiceGenerationStrategy, error) {
	if invoiceType == InvoiceTypeCash {
		return CashBasisInvoiceGeneration{}, nil
	}

	if invoiceType == InvoiceTypeAccrual {
		return AccrualInvoiceGeneration{}, nil
	}

	return nil, fmt.Errorf("Invoice Type %s is invalid", invoiceType)
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

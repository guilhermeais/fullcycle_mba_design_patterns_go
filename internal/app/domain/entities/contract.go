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
	payments        []Payment
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

func (c Contract) GetBalance() float64 {
	balance := c.Amount

	for _, p := range c.payments {
		balance -= p.Amount
	}

	return balance
}

func (c *Contract) AddPayment(p Payment) {
	c.payments = append(c.payments, p)
}

func (c Contract) GetPayments() []Payment {
	return c.payments
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

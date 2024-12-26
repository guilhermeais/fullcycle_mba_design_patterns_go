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

func (c Contract) GenerateInvoices(month, year int, invoiceType InvoiceType) (invoices []Invoice) {
	if invoiceType == InvoiceTypeCash {
		for _, payment := range c.Payments {
			if int(payment.Date.Month()) != month || payment.Date.Year() != year {
				continue
			}
			invoices = append(invoices, Invoice{Date: payment.Date, Amount: payment.Amount})
		}
	}

	if invoiceType == InvoiceTypeAccrual {
		period := 0
		for period <= c.Periods {
			date := c.Date.AddDate(0, period, 0)
			period++
			amount := c.Amount / float64(c.Periods)

			if int(date.Month()) != month || date.Year() != year {
				continue
			}

			invoices = append(invoices, Invoice{Date: date, Amount: amount})
		}
	}

	return invoices
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

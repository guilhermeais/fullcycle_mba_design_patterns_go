package usecase

import (
	"context"
	domain "invoices/internal/app/domain/entities"
	"time"
)

type GenerateInvoicesInput struct {
	Year  int                `json:"year"`
	Month int                `json:"month"`
	Type  domain.InvoiceType `json:"type"`
}

type GenerateInvoicesOutput struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type GenerateInvoices struct {
	contractRepository domain.ContractRepository
	observer           *Observer[InvoiceGeneratedEventData]
}

const resultDateFormat = "2006-01-02"

func (gi *GenerateInvoices) Execute(input GenerateInvoicesInput) ([]GenerateInvoicesOutput, error) {
	contracts, err := gi.contractRepository.List(context.Background())
	if err != nil {
		return nil, err
	}
	results := []GenerateInvoicesOutput{}
	for _, c := range contracts {
		invoices, err := c.GenerateInvoices(input.Month, input.Year, input.Type)
		if err != nil {
			return nil, err
		}
		for _, invoice := range invoices {
			go gi.observer.Notify(Event[InvoiceGeneratedEventData]{
				Type: InvoiceGenerated,
				Date: time.Now(),
				Data: InvoiceGeneratedEventData{
					Amount:    invoice.Amount,
					Date:      invoice.Date,
					UserEmail: "guilhermeteixeiraais@gmail.com",
				},
			})

			results = append(results, GenerateInvoicesOutput{Date: invoice.Date.Format(resultDateFormat), Amount: invoice.Amount})
		}
	}

	return results, nil
}

func NewGenerateInvoices(repo domain.ContractRepository, observer *Observer[InvoiceGeneratedEventData]) *GenerateInvoices {
	return &GenerateInvoices{repo, observer}
}

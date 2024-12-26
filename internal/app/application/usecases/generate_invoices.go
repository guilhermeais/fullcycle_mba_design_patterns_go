package usecase

import (
	"context"
	domain "invoices/internal/app/domain/entities"
)

type GenerateInvoicesInput struct {
	Year, Month int
	Type        domain.InvoiceType
}

type GenerateInvoicesOutput struct {
	Date   string
	Amount float64
}

type GenerateInvoices struct {
	contractRepository domain.ContractRepository
}

const resultDateFormat = "2006-01-02"

func (generateInvioices *GenerateInvoices) Execute(input GenerateInvoicesInput) ([]GenerateInvoicesOutput, error) {
	contracts, err := generateInvioices.contractRepository.List(context.Background())
	if err != nil {
		return nil, err
	}
	var results []GenerateInvoicesOutput
	for _, c := range contracts {
		invoices, err := c.GenerateInvoices(input.Month, input.Year, input.Type)
		if err != nil {
			return nil, err
		}
		for _, invoice := range invoices {
			results = append(results, GenerateInvoicesOutput{Date: invoice.Date.Format(resultDateFormat), Amount: invoice.Amount})
		}
	}

	return results, nil
}

func NewGenerateInvoices(repo domain.ContractRepository) *GenerateInvoices {
	return &GenerateInvoices{repo}
}

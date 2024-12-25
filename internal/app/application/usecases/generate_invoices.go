package usecase

import (
	"context"
	domain "invoices/internal/app/domain/entities"
)

type GenerateInvoicesInput struct {
	Year, Month int
	Type        string
}

type GenerateInvoicesOutput struct {
	Date   string
	Amount float64
}

type GenerateInvoices struct {
	contractRepository domain.ContractRepository
}

func (generateInvioices *GenerateInvoices) Execute(input GenerateInvoicesInput) ([]GenerateInvoicesOutput, error) {
	contracts, err := generateInvioices.contractRepository.List(context.Background())
	if err != nil {
		return nil, err
	}
	var results []GenerateInvoicesOutput
	for _, c := range contracts {
		const resultDateFormat = "2006-01-02"

		for _, invoice := range c.GenerateInvoices(input.Month, input.Year, input.Type) {
			results = append(results, GenerateInvoicesOutput{Date: invoice.Date.Format(resultDateFormat), Amount: invoice.Amount})
		}
	}

	return results, nil
}

func NewGenerateInvoices(repo domain.ContractRepository) *GenerateInvoices {
	return &GenerateInvoices{repo}
}

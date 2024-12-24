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
	Date  string
	Amout float64
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
	for _, contract := range contracts {
		const resultDateFormat = "2006-01-02"

		for _, payment := range contract.Payments {
			if input.Type == "cash" {
				if int(payment.Date.Month()) != input.Month || payment.Date.Year() != input.Year {
					continue
				}
				results = append(results, GenerateInvoicesOutput{Date: payment.Date.Format(resultDateFormat), Amout: payment.Amount})
			}

			if input.Type == "accrual" {
				period := 0

				for period <= contract.Periods {
					date := contract.Date.AddDate(0, period, 0)
					period++
					amount := contract.Amount / float64(contract.Periods)

					if int(date.Month()) != input.Month || date.Year() != input.Year {
						continue
					}

					results = append(results, GenerateInvoicesOutput{Date: date.Format(resultDateFormat), Amout: amount})
				}
			}
		}
	}

	return results, nil
}

func NewGenerateInvoices(repo domain.ContractRepository) *GenerateInvoices {
	return &GenerateInvoices{repo}
}

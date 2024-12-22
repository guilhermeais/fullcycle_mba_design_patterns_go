package usecase

import (
	"context"
	"fmt"
	domain "invoices/internal/app/domain/entities"
	"os"

	"github.com/jackc/pgx/v5"
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
}

const (
	getContractsQuery = "select * from invoices_service.contract"
	paymentQuery      = "select id, date, amount from invoices_service.payment where contract_id = $1"
)

func (uc *GenerateInvoices) Execute(input GenerateInvoicesInput) ([]GenerateInvoicesOutput, error) {
	dbUrl := os.Getenv("POSTGRES_URL")
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %s: %v", dbUrl, err)
	}
	defer conn.Close(context.Background())

	var results []GenerateInvoicesOutput

	var contracts []domain.Contract
	contractRows, err := conn.Query(context.Background(), getContractsQuery)
	if err != nil {
		return nil, fmt.Errorf("unable to get contracts: %v", err)
	}
	defer contractRows.Close()
	for contractRows.Next() {
		var contract domain.Contract
		err := contractRows.Scan(&contract.Id, &contract.Description, &contract.Amount, &contract.Periods, &contract.Date)
		if err != nil {
			return nil, fmt.Errorf("unable to scan contract: %v", err)
		}
		contracts = append(contracts, contract)
	}

	if contractRows.Err() != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", contractRows.Err())
	}

	for _, contract := range contracts {
		paymentRows, err := conn.Query(context.Background(), paymentQuery, contract.Id)
		if err != nil {
			return nil, fmt.Errorf("unable to get payments: %v", err)
		}
		defer paymentRows.Close()

		const resultDateFormat = "2006-01-02"

		for paymentRows.Next() {
			if input.Type == "cash" {
				var payment domain.Payment
				err := paymentRows.Scan(&payment.Id, &payment.Date, &payment.Amount)
				if err != nil {
					return nil, fmt.Errorf("unable to scan payment: %v", err)
				}

				fmt.Printf("payment date: %v", int(payment.Date.Month()))

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

func NewGenerateInvoices() *GenerateInvoices {
	return &GenerateInvoices{}
}

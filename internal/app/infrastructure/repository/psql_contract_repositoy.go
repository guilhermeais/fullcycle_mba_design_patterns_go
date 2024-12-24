package repository

import (
	"context"
	"fmt"
	domain "invoices/internal/app/domain/entities"
	"os"

	"github.com/jackc/pgx/v5"
)

type PSQLContractRepository struct {
}

const (
	getContractsQuery        = "select * from invoices_service.contract"
	getContractPaymentsQuery = "select id, date, amount from invoices_service.payment where contract_id = $1"
)

func (PSQLContractRepository) List(ctx context.Context) ([]domain.Contract, error) {
	dbUrl := os.Getenv("POSTGRES_URL")
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %s: %v", dbUrl, err)
	}
	defer conn.Close(ctx)
	var contracts []domain.Contract
	contractRows, err := conn.Query(ctx, getContractsQuery)
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

	for i, contract := range contracts {
		paymentRows, err := conn.Query(context.Background(), getContractPaymentsQuery, contract.Id)
		if err != nil {
			return nil, fmt.Errorf("unable to get payments: %v", err)
		}
		defer paymentRows.Close()

		for paymentRows.Next() {
			var payment domain.Payment
			err := paymentRows.Scan(&payment.Id, &payment.Date, &payment.Amount)
			if err != nil {
				return nil, fmt.Errorf("unable to scan payment: %v", err)
			}

			contracts[i].Payments = append(contract.Payments, payment)
		}
	}

	return contracts, nil
}

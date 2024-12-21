package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type GenerateInvoices struct {
}

// id uuid not null default uuid_generate_v4() primary key,
// description text,
// amount numeric,
// periods integer,
// date timestamp

type Contract struct {
	id, description string
	amount          float64
	periods         int
	date            time.Time
}

func (uc *GenerateInvoices) Execute() ([]Contract, error) {
	dbUrl := os.Getenv("POSTGRES_URL")
	fmt.Printf("connecting to '%s'", dbUrl)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %s: %v", dbUrl, err)
	}
	defer conn.Close(context.Background())

	const query = "select * from invoices_service.contract"
	contracts := []Contract{}
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("unable to get contracts: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var contract Contract
		err := rows.Scan(&contract.id, &contract.description, &contract.amount, &contract.periods, &contract.date)
		if err != nil {
			return nil, fmt.Errorf("unable to scan contract: %v", err)
		}
		contracts = append(contracts, contract)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", rows.Err())
	}

	return contracts, nil
}

func NewGenerateInvoices() *GenerateInvoices {
	return &GenerateInvoices{}
}

package testutils

import (
	"context"
	"fmt"

	domain "invoices/internal/app/domain/entities"

	"github.com/jackc/pgx/v5"
)

const createContractSQL = `
    INSERT INTO contract (id, description, amount, periods, date)
    VALUES ($1, $2, $3, $4, $5);
`

const createPaymentSQL = `
    INSERT INTO payment (id, contract_id, amount, date)
    VALUES ($1, $2, $3, $4);
`

type ContractFactory struct {
	Conn *pgx.Conn
}

func (cf *ContractFactory) CreateContract(c domain.Contract) error {
	ctx := context.Background()
	tx, err := cf.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaciton: %v", err)
	}
	_, err = tx.Exec(ctx, createContractSQL, c.Id, c.Description, c.Amount, c.Periods, c.Date)
	if err != nil {
		return fmt.Errorf("error creating contract: %v", err)
	}

	for _, p := range c.GetPayments() {
		_, err = tx.Exec(ctx, createPaymentSQL, p.Id, c.Id, p.Amount, p.Date)
		if err != nil {
			return fmt.Errorf("error creating payment: %v", err)
		}
	}

	tx.Commit(ctx)
	return nil
}

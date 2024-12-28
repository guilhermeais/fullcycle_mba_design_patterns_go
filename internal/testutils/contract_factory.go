package testutils

import (
	"context"
	"fmt"

	domain "invoices/internal/app/domain/entities"
	"invoices/internal/app/infrastructure/repository"
)

const createContractSQL = `
    INSERT INTO invoices_service.contract (id, description, amount, periods, date)
    VALUES ($1, $2, $3, $4, $5);
`

const createPaymentSQL = `
    INSERT INTO invoices_service.payment (id, contract_id, amount, date)
    VALUES ($1, $2, $3, $4);
`

type ContractFactory struct{}

func (ContractFactory) CreateContract(c domain.Contract) error {
	ctx := context.Background()
	conn, err := repository.MakePGConnectionWithUri(PgContainer.URI)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
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

package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	usecases "invoices/internal/app/application/usecases"
	domain "invoices/internal/app/domain/entities"
)

type InMemoryContractRepository struct {
	Contracts []domain.Contract
}

func (r InMemoryContractRepository) List(ctx context.Context) ([]domain.Contract, error) {
	return r.Contracts, nil
}

func TestGenerateInvoices(t *testing.T) {
	t.Run("Deve gerar notas fiscais por regime de caixa", func(t *testing.T) {
		mockedContract := domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		}
		mockedContract.AddPayment(domain.Payment{
			Id:     "6355b223-fce0-4f7c-998a-1f027281e308",
			Amount: 6000,
			Date:   time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		})
		generateInvoices, _ := makeSut(t, mockedContract)
		input := usecases.GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  domain.InvoiceTypeCash,
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2024-12-18", output[0].Date)
		assert.Equal(t, float64(6000), output[0].Amount)
	})

	t.Run("Deve gerar notas fiscais por regime de competência", func(t *testing.T) {
		generateInvoices, _ := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		})

		input := usecases.GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  domain.InvoiceTypeAccrual,
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2024-12-18", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amount)
	})

	t.Run("Deve gerar notas fiscais por regime de competência (última data)", func(t *testing.T) {
		generateInvoices, _ := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		})
		input := usecases.GenerateInvoicesInput{
			Year:  2025,
			Month: 12,
			Type:  domain.InvoiceTypeAccrual,
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2025-12-18", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amount)
	})

	t.Run("Deve gerar notas fiscais por regime de competência (fora do periodo)", func(t *testing.T) {
		generateInvoices, _ := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		})
		input := usecases.GenerateInvoicesInput{
			Year:  2026,
			Month: 1,
			Type:  domain.InvoiceTypeAccrual,
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 0)
	})

	t.Run("Deve notificar os inscritos em InvoiceGenerated quando a nota fiscal for gerada por regime de competência", func(t *testing.T) {
		generateInvoices, observer := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
		})

		invoiceChan := make(chan usecases.Event[usecases.InvoiceGeneratedEventData], 1)
		observer.Subscribe(usecases.InvoiceGenerated, invoiceChan)

		input := usecases.GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  domain.InvoiceTypeAccrual,
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2024-12-18", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amount)

		invoiceEvent := <-invoiceChan
		assert.Equal(t, usecases.InvoiceGenerated, invoiceEvent.Type)
		assert.Equal(t, time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC), invoiceEvent.Data.Date)
		assert.Equal(t, float64(500), invoiceEvent.Data.Amount)
		assert.Equal(t, "guilhermeteixeiraais@gmail.com", invoiceEvent.Data.UserEmail)
	})
}

func makeSut(t *testing.T, mockedContracts ...domain.Contract) (*usecases.GenerateInvoices, *usecases.Observer[usecases.InvoiceGeneratedEventData]) {
	t.Helper()
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	observer := usecases.NewObserver[usecases.InvoiceGeneratedEventData]()
	generateInvoices := usecases.NewGenerateInvoices(InMemoryContractRepository{
		Contracts: mockedContracts,
	}, observer)

	return generateInvoices, observer
}

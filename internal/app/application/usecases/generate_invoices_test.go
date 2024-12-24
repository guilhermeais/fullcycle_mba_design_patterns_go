package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

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
		generateInvoices := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			Payments: []domain.Payment{{
				Id:     "6355b223-fce0-4f7c-998a-1f027281e308",
				Amount: 6000,
				Date:   time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			}},
		})
		input := GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  "cash",
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2024-12-18", output[0].Date)
		assert.Equal(t, float64(6000), output[0].Amout)
	})

	t.Run("Deve gerar notas fiscais por regime de competência", func(t *testing.T) {
		generateInvoices := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			Payments: []domain.Payment{{
				Id:     "6355b223-fce0-4f7c-998a-1f027281e308",
				Amount: 6000,
				Date:   time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			}},
		})
		input := GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  "accrual",
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2024-12-18", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amout)
	})

	t.Run("Deve gerar notas fiscais por regime de competência (última data)", func(t *testing.T) {
		generateInvoices := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			Payments: []domain.Payment{{
				Id:     "6355b223-fce0-4f7c-998a-1f027281e308",
				Amount: 6000,
				Date:   time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			}},
		})
		input := GenerateInvoicesInput{
			Year:  2025,
			Month: 12,
			Type:  "accrual",
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2025-12-18", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amout)
	})

	t.Run("Deve gerar notas fiscais por regime de competência (fora do periodo)", func(t *testing.T) {
		generateInvoices := makeSut(t, domain.Contract{
			Id:          "fac05a57-7d61-4283-ab32-7696902eac44",
			Description: "prestação de serviços escolares",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			Payments: []domain.Payment{{
				Id:     "6355b223-fce0-4f7c-998a-1f027281e308",
				Amount: 6000,
				Date:   time.Date(2024, time.Month(12), 18, 10, 0, 0, 0, time.UTC),
			}},
		})
		input := GenerateInvoicesInput{
			Year:  2026,
			Month: 1,
			Type:  "accrual",
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 0)
	})
}

func makeSut(t *testing.T, mockedContracts ...domain.Contract) *GenerateInvoices {
	t.Helper()
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	generateInvoices := NewGenerateInvoices(InMemoryContractRepository{
		Contracts: mockedContracts,
	})

	return generateInvoices
}

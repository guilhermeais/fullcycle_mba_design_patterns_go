package domain_test

import (
	"testing"
	"time"

	domain "invoices/internal/app/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestGenerateInvoices(t *testing.T) {
	t.Run("Deve gerar faturas de um contrato", func(t *testing.T) {
		contract := domain.Contract{
			Id:          "",
			Description: "",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, 12, 25, 17, 10, 0, 0, time.UTC),
		}

		firstInvoice, _ := contract.GenerateInvoices(12, 2024, domain.InvoiceTypeAccrual)
		secondInvoice, _ := contract.GenerateInvoices(1, 2025, domain.InvoiceTypeAccrual)

		assert.Equal(t, time.Date(2024, 12, 25, 17, 10, 0, 0, time.UTC), firstInvoice[0].Date)
		assert.Equal(t, time.Date(2025, 1, 25, 17, 10, 0, 0, time.UTC), secondInvoice[0].Date)
	})

	t.Run("Deve calcular o saldo do contrato", func(t *testing.T) {
		contract := domain.Contract{
			Id:          "",
			Description: "",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, 12, 25, 17, 10, 0, 0, time.UTC),
		}

		assert.Equal(t, float64(6000), contract.GetBalance())

		contract.AddPayment(domain.Payment{
			Id:     "",
			Amount: 2000,
			Date:   time.Date(2024, 12, 25, 17, 10, 0, 0, time.UTC),
		})

		assert.Equal(t, float64(4000), contract.GetBalance())
	})
}

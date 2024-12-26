package domain_test

import (
	"testing"
	"time"

	domain "invoices/internal/app/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestGenerateInvoices(t *testing.T) {
	t.Run("Deve gerar faturas de um contrat", func(t *testing.T) {
		contract := domain.Contract{
			Id:          "",
			Description: "",
			Amount:      6000,
			Periods:     12,
			Date:        time.Date(2024, 12, 25, 17, 10, 0, 0, time.UTC),
		}

		assert.Equal(t, time.Date(2024, 12, 25, 17, 10, 0, 0, time.UTC), contract.GenerateInvoices(12, 2024, domain.InvoiceTypeAccrual)[0].Date)
		assert.Equal(t, time.Date(2025, 1, 25, 17, 10, 0, 0, time.UTC), contract.GenerateInvoices(1, 2025, domain.InvoiceTypeAccrual)[0].Date)
	})
}

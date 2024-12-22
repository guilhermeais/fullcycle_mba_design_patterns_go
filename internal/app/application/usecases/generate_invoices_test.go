package usecase

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateInvoices(t *testing.T) {
	t.Run("Deve gerar notas fiscais por regime de caixa", func(t *testing.T) {
		err := godotenv.Load("../../../../.env")
		if err != nil {
			t.Fatalf("Error loading .env file: %v", err)
		}
		generateInvoices := NewGenerateInvoices()
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
		err := godotenv.Load("../../../../.env")
		if err != nil {
			t.Fatalf("Error loading .env file: %v", err)
		}
		generateInvoices := NewGenerateInvoices()
		input := GenerateInvoicesInput{
			Year:  2024,
			Month: 12,
			Type:  "accrual",
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2024-12-19", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amout)
	})

	t.Run("Deve gerar notas fiscais por regime de competência (última data)", func(t *testing.T) {
		err := godotenv.Load("../../../../.env")
		if err != nil {
			t.Fatalf("Error loading .env file: %v", err)
		}
		generateInvoices := NewGenerateInvoices()
		input := GenerateInvoicesInput{
			Year:  2025,
			Month: 12,
			Type:  "accrual",
		}
		output, err := generateInvoices.Execute(input)

		assert.Nil(t, err, fmt.Sprintf("does not expected the error: %v", err))
		assert.Len(t, output, 1)

		assert.Equal(t, "2025-12-19", output[0].Date)
		assert.Equal(t, float64(500), output[0].Amout)
	})

	t.Run("Deve gerar notas fiscais por regime de competência (fora do periodo)", func(t *testing.T) {
		err := godotenv.Load("../../../../.env")
		if err != nil {
			t.Fatalf("Error loading .env file: %v", err)
		}
		generateInvoices := NewGenerateInvoices()
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

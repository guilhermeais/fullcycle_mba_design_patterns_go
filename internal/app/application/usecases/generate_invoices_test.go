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

		assert.Equal(t, output[0].Date, "2024-12-18")
		assert.Equal(t, output[0].Amout, 6000)
	})
}

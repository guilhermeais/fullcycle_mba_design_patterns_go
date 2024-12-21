package usecase

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGenerateInvoices(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	generateInvoices := NewGenerateInvoices()
	output, err := generateInvoices.Execute()

	if err != nil {
		t.Fatalf("does not expected the error: %v", err)
	}

	if len(output) <= 0 {
		t.Fatal("expect output to have length of 0")
	}

}

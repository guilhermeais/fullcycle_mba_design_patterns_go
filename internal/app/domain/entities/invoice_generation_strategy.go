package domain

type InvoiceGenerationStrategy interface {
	Generate(c Contract, month, year int) []Invoice
}

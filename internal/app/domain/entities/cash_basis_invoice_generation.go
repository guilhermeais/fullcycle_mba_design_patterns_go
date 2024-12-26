package domain

type CashBasisInvoiceGeneration struct{}

func (CashBasisInvoiceGeneration) Generate(c Contract, month, year int) (invoices []Invoice) {
	for _, payment := range c.Payments {
		if int(payment.Date.Month()) != month || payment.Date.Year() != year {
			continue
		}
		invoices = append(invoices, Invoice{Date: payment.Date, Amount: payment.Amount})
	}

	return invoices
}

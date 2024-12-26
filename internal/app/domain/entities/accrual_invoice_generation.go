package domain

type AccrualInvoiceGeneration struct{}

func (AccrualInvoiceGeneration) Generate(c Contract, month, year int) (invoices []Invoice) {
	period := 0
	for period <= c.Periods {
		date := c.Date.AddDate(0, period, 0)
		period++
		amount := c.Amount / float64(c.Periods)

		if int(date.Month()) != month || date.Year() != year {
			continue
		}

		invoices = append(invoices, Invoice{Date: date, Amount: amount})
	}

	return invoices
}

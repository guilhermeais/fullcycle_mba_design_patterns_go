package domain

import "fmt"

func MakeInvoiceGenerationStrategy(invoiceType InvoiceType) (InvoiceGenerationStrategy, error) {
	if invoiceType == InvoiceTypeCash {
		return CashBasisInvoiceGeneration{}, nil
	}

	if invoiceType == InvoiceTypeAccrual {
		return AccrualInvoiceGeneration{}, nil
	}

	return nil, fmt.Errorf("Invoice Type %s is invalid", invoiceType)
}

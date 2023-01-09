package main

import (
	"invoicedata/model"
	"io"
)

type InvoicesMarshaler interface {
	MarshalInvoices(writer io.Writer, invoices []*model.Invoice)
}
type InvoiceUnmarshaler interface {
	UnmarshalInvoices(reader io.Reader) ([]*model.Invoice, error)
}

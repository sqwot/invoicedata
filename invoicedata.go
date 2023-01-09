package main

import (
	"fmt"
	"invoicedata/model"
	"io"
)

func readInvoices(reader io.Reader, suffix string) ([]*model.Invoice, error) {
	var unmarshaler InvoiceUnmarshaler
	switch suffix {
	case ".gob":
		unmarshaler = GobMarshaler{}
	case ".inv":
		unmarshaler = InvMarshaler{}
	case ".jsn", ".json":
		unmarshaler = JSONMarshaler{}
	case ".txt":
		unmarshaler = TXTMarshaler{}
	case ".xml":
		unmarshaler = XMLMarshaler{}
	}
	if unmarshaler != nil {
		return unmarshaler.UnmarshalInvoices(reader)
	}
	return nil, fmt.Errorf("unrecognized input suffix: %s", suffix)
}

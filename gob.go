package main

import (
	"encoding/gob"
	"fmt"
	"invoicedata/constants"
	"invoicedata/model"
	"io"
)

type GobMarshaler struct{}

func (GobMarshaler) MarshalInvoices(writer io.Writer, invoices []*model.Invoice) error {
	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(constants.MagicNumber); err != nil {
		return err
	}
	if err := encoder.Encode(constants.FileVersion); err != nil {
		return err
	}

	return encoder.Encode(invoices)
}

func (GobMarshaler) UnmarshalInvoices(reader io.Reader) ([]*model.Invoice, error) {
	decoder := gob.NewDecoder(reader)
	var magic int
	if err := decoder.Decode(&magic); err != nil {
		return nil, err
	}
	if magic != constants.MagicNumber {
		return nil, fmt.Errorf("cannot read non-invoices gob file")
	}
	var version int
	if err := decoder.Decode(&version); err != nil {
		return nil, err
	}
	if version > constants.FileVersion {
		return nil, fmt.Errorf("version %d is too new to read", version)
	}
	var invoices []*model.Invoice
	err := decoder.Decode(&invoices)
	return invoices, err
}

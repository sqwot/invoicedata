package main

import (
	"encoding/json"
	"fmt"
	"invoicedata/constants"
	"invoicedata/model"
	"io"
)

type JSONMarshaler struct{}

func (JSONMarshaler) MarshalInvoices(writer io.Writer, invoices []*model.Invoice) error {
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(constants.FileType); err != nil {
		return err
	}
	if err := encoder.Encode(constants.FileVersion); err != nil {
		return err
	}
	return encoder.Encode(invoices)
}

func (JSONMarshaler) UnmarshalInvoices(reader io.Reader) ([]*model.Invoice, error) {
	decoder := json.NewDecoder(reader)
	var kind string
	if err := decoder.Decode(&kind); err != nil {
		return nil, err
	}
	if kind != constants.FileType {
		return nil, fmt.Errorf("cannot read non-invoices json file")
	}
	var version int
	if err := decoder.Decode(&version); err != nil {
		return nil, err
	}
	if version > constants.FileVersion {
		return nil, fmt.Errorf("version %d is too new for read", version)
	}
	var invoices []*model.Invoice
	err := decoder.Decode(&invoices)
	return invoices, err
}

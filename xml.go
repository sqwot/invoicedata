package main

import (
	"encoding/xml"
	"fmt"
	"invoicedata/constants"
	"invoicedata/model"
	"io"
)

type XMLMarshaler struct{}

func (XMLMarshaler) MarshalInvoices(writer io.Writer, invoices []*model.Invoice) error {
	if _, err := writer.Write([]byte(xml.Header)); err != nil {
		return err
	}
	xmlInvoices := model.XMLInvoicesForInvoices(invoices)
	encoder := xml.NewEncoder(writer)
	return encoder.Encode(xmlInvoices)
}

func (XMLMarshaler) UnmarshalInvoices(reader io.Reader) ([]*model.Invoice, error) {
	xmlInvoices := &model.XMLInvoices{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(xmlInvoices); err != nil {
		return nil, err
	}
	if xmlInvoices.Version > constants.FileVersion {
		return nil, fmt.Errorf("version %d is too new to read", xmlInvoices.Version)
	}
	return xmlInvoices.Invoices()
}

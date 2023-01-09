package model

import (
	"encoding/json"
	"encoding/xml"
	"invoicedata/constants"
	"time"
)

type Invoice struct {
	Id         int
	CustomerId int
	Raised     time.Time
	Due        time.Time
	Paid       bool
	Note       string
	Items      []*Item
}

type JSONInvoice struct {
	Id         int
	CustomerId int
	Raised     string
	Due        string
	Paid       bool
	Note       string
	Items      []*Item
}

func (invoice Invoice) MarshalJSON() ([]byte, error) {
	jsonInvoice := JSONInvoice{
		invoice.Id,
		invoice.CustomerId,
		invoice.Raised.Format(constants.DateFormat),
		invoice.Due.Format(constants.DateFormat),
		invoice.Paid,
		invoice.Note,
		invoice.Items,
	}
	return json.Marshal(jsonInvoice)
}

func (invoice *Invoice) UnmarshalJSON(data []byte) (err error) {
	var jsonInvoice JSONInvoice
	if err := json.Unmarshal(data, &jsonInvoice); err != nil {
		return err
	}
	var raised, due time.Time
	if raised, err = time.Parse(constants.DateFormat, jsonInvoice.Raised); err != nil {
		return err
	}
	if due, err = time.Parse(constants.DateFormat, jsonInvoice.Due); err != nil {
		return err
	}
	*invoice = Invoice{
		jsonInvoice.Id,
		jsonInvoice.CustomerId,
		raised,
		due,
		jsonInvoice.Paid,
		jsonInvoice.Note,
		jsonInvoice.Items,
	}
	return nil
}

type XMLInvoices struct {
	XMLName xml.Name      `xml:"INVOICES"`
	Version int           `xml:"version,attr"`
	Invoice []*XMLInvoice `xml:"INVOICE"`
}

type XMLInvoice struct {
	XMLName    xml.Name   `xml:"INVOICE"`
	Id         int        `xml:",attr"`
	CustomerId int        `xml:",attr"`
	Raised     string     `xml:",attr"`
	Due        string     `xml:",attr"`
	Paid       bool       `xml:",attr"`
	Note       string     `xml:"NOTE"`
	Item       []*XMLItem `xml:"ITEM"`
}

func XMLInvoicesForInvoices(invoices []*Invoice) *XMLInvoices {
	xmlInvoices := &XMLInvoices{
		Version: constants.FileVersion,
		Invoice: make([]*XMLInvoice, 0, len(invoices)),
	}
	for _, invoice := range invoices {
		xmlInvoices.Invoice = append(xmlInvoices.Invoice, XMLInvoiceForInvoice(invoice))
	}
	return xmlInvoices
}

func XMLInvoiceForInvoice(invoice *Invoice) *XMLInvoice {
	xmlInvoice := &XMLInvoice{
		Id:         invoice.Id,
		CustomerId: invoice.CustomerId,
		Raised:     invoice.Raised.Format(constants.DateFormat),
		Due:        invoice.Due.Format(constants.DateFormat),
		Paid:       invoice.Paid,
		Note:       invoice.Note,
		Item:       make([]*XMLItem, 0, len(invoice.Items)),
	}
	for _, item := range invoice.Items {
		xmlItem := &XMLItem{
			Id:       item.Id,
			Price:    item.Price,
			Quantity: item.Quantity,
			Note:     item.Note,
		}
		xmlInvoice.Item = append(xmlInvoice.Item, xmlItem)
	}
	return xmlInvoice
}

func (xmlInvoices *XMLInvoices) Invoices() (invoices []*Invoice, err error) {
	invoices = make([]*Invoice, 0, len(xmlInvoices.Invoice))
	for _, xmlInvoice := range xmlInvoices.Invoice {
		invoice, err := xmlInvoice.Invoice()
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func (xmlInvoice *XMLInvoice) Invoice() (invoice *Invoice, err error) {
	invoice = &Invoice{
		Id:         xmlInvoice.Id,
		CustomerId: xmlInvoice.CustomerId,
		Paid:       xmlInvoice.Paid,
		Note:       xmlInvoice.Note,
		Items:      make([]*Item, 0, len(xmlInvoice.Item)),
	}
	if invoice.Raised, err = time.Parse(constants.DateFormat, xmlInvoice.Raised); err != nil {
		return nil, err
	}
	if invoice.Due, err = time.Parse(constants.DateFormat, xmlInvoice.Due); err != nil {
		return nil, err
	}
	for _, xmlItem := range xmlInvoice.Item {
		item := &Item{
			Id:       xmlItem.Id,
			Price:    xmlItem.Price,
			Quantity: xmlItem.Quantity,
			Note:     xmlItem.Note,
		}
		invoice.Items = append(invoice.Items, item)
	}
	return invoice, nil
}

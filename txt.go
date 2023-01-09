package main

import (
	"bufio"
	"fmt"
	"invoicedata/constants"
	"invoicedata/model"
	"io"
	"strings"
	"time"
)

const (
	noteSep = ":"
)

type TXTMarshaler struct{}
type writerFunc func(string, ...interface{}) error

func (TXTMarshaler) MarshalInvoices(writer io.Writer, invoices []*model.Invoice) error {
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()
	var write writerFunc = func(format string, args ...interface{}) error {
		_, err := fmt.Fprintf(bufferedWriter, format, args...)
		return err
	}
	if err := write("%s %d\n", constants.FileType, constants.FileVersion); err != nil {
		return err
	}
	for _, invoice := range invoices {
		if err := write.writeInvoice(invoice); err != nil {
			return err
		}
	}
	return nil
}

func (TXTMarshaler) UnmarshalInvoices(reader io.Reader) ([]*model.Invoice, error) {
	bufferedReader := bufio.NewReader(reader)
	if err := checkTxtVersion(bufferedReader); err != nil {
		return nil, err
	}
	var invoices []*model.Invoice
	eof := false
	for lino := 2; !eof; lino++ {
		line, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return nil, err
		}
		if invoices, err = parceTxtLine(lino, line, invoices); err != nil {
			return nil, err
		}
	}
	return invoices, nil
}

func (write writerFunc) writeInvoice(invoice *model.Invoice) error {
	note := ""
	if invoice.Note != "" {
		note = noteSep + " " + invoice.Note
	}
	if err := write("INVOICE ID=%d CUSTOMER=%d RAISED=%s DUE=%s PAID=%t%s\n",
		invoice.Id, invoice.CustomerId,
		invoice.Raised.Format(constants.DateFormat),
		invoice.Due.Format(constants.DateFormat),
		invoice.Paid, note); err != nil {
		return err
	}

	if err := write.writeItems(invoice.Items); err != nil {
		return err
	}
	return write("\f\n")
}

func (write writerFunc) writeItems(items []*model.Item) error {
	for _, item := range items {
		note := ""
		if item.Note != "" {
			note = noteSep + " " + note
		}

		if err := write("ITEM ID=%s PRICE=%.2f QUANTITY=%d%s\n",
			item.Id, item.Price, item.Quantity, note); err != nil {
			return err
		}
	}
	return nil
}

func checkTxtVersion(reader io.Reader) error {
	var version int
	if _, err := fmt.Fscanf(reader, "INVOICES %d\n", &version); err != nil {
		return fmt.Errorf("cannot read non-invioces text file %s", err.Error())
	} else if version > constants.FileVersion {
		return fmt.Errorf("version %d is too new to read", version)
	}
	return nil

}

func parceTxtLine(lino int, line string, invoices []*model.Invoice) ([]*model.Invoice, error) {
	var err error

	if strings.HasPrefix(line, "INVOICE") {
		var invoice *model.Invoice
		invoice, err = parceTxtInvoice(lino, line)
		invoices = append(invoices, invoice)
	} else if strings.HasPrefix(line, "ITEM") {
		if len(invoices) == 0 {
			err = fmt.Errorf("item outside of an invoice line %d", lino)
		} else {
			var item *model.Item
			item, err = parseTxtItem(lino, line)
			items := &invoices[len(invoices)-1].Items
			*items = append(*items, item)
		}
	}
	return invoices, err
}

func parceTxtInvoice(lino int, line string) (invoice *model.Invoice, err error) {
	invoice = &model.Invoice{}
	var raised, due string
	if _, err = fmt.Sscanf(line, "INVOICE ID=&d CUSTOMER=%d"+
		"RAISED=%s DUE=%s PAID=%t ",
		&invoice.Id, &invoice.CustomerId, &raised, &due, &invoice.Paid); err != nil {
		return nil, fmt.Errorf("invalid invoice %v line %d", err, lino)
	}

	if invoice.Raised, err = time.Parse(constants.DateFormat, raised); err != nil {
		return nil, fmt.Errorf("invalid raised %v line %d", err, lino)
	}

	if invoice.Due, err = time.Parse(constants.DateFormat, due); err != nil {
		return nil, fmt.Errorf("invalid due %v line %d", err, lino)
	}

	if i := strings.Index(line, noteSep); i > -1 {
		invoice.Note = strings.TrimSpace(line[i+len(noteSep):])
	}

	return invoice, nil
}

func parseTxtItem(lino int, line string) (item *model.Item, err error) {
	item = &model.Item{}
	if _, err = fmt.Sscanf(line, "ITEM ID=%s PRICE=%f QUANTITY=%d", &item.Id, &item.Price, &item.Quantity); err != nil {
		return nil, fmt.Errorf("invalid item %v line %d", err, lino)
	}
	if i := strings.Index(line, noteSep); i > -1 {
		item.Note = strings.TrimSpace(line[i+len(noteSep):])
	}
	return item, nil
}

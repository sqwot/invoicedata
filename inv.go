package main

import (
	"encoding/binary"
	"fmt"
	"invoicedata/constants"
	"invoicedata/model"
	"io"
	"strconv"
	"time"
)

const (
	invDateFormat = "20060102"
)

var byteOrder = binary.LittleEndian

type InvMarshaler struct{}
type invWriteFunc func(interface{}) error

func (InvMarshaler) MarshalInvoices(writer io.Writer, invoices []*model.Invoice) error {
	var write invWriteFunc = func(x interface{}) error {
		return binary.Write(writer, byteOrder, x)
	}

	if err := write(uint32(constants.MagicNumber)); err != nil {
		return err
	}

	if err := write(uint16(constants.FileVersion)); err != nil {
		return err
	}

	if err := write(uint32(len(invoices))); err != nil {
		return err
	}
	for _, invoice := range invoices {
		if err := write.writeInvoice(invoice); err != nil {
			return err
		}
	}
	return nil
}

func (write invWriteFunc) writeInvoice(invoice *model.Invoice) error {
	for _, i := range []int{invoice.Id, invoice.CustomerId} {
		if err := write(int32(i)); err != nil {
			return err
		}
	}
	for _, date := range []time.Time{invoice.Raised, invoice.Due} {
		if err := write.writeDate(date); err != nil {
			return err
		}
	}
	if err := write.writeBool(invoice.Paid); err != nil {
		return err
	}
	if err := write.writeString(invoice.Note); err != nil {
		return err
	}
	if err := write(int32(len(invoice.Items))); err != nil {
		return err
	}
	for _, item := range invoice.Items {
		if err := write.writeItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (write invWriteFunc) writeDate(date time.Time) error {
	i, err := strconv.Atoi(date.Format(invDateFormat))
	if err != nil {
		return err
	}
	return write(int32(i))
}

func (write invWriteFunc) writeBool(b bool) error {
	var v int8
	if b {
		v = 1
	}
	return write(v)
}

func (write invWriteFunc) writeString(s string) error {
	if err := write(int32(len(s))); err != nil {
		return err
	}
	return write([]byte(s))
}

func (write invWriteFunc) writeItem(item *model.Item) error {
	if err := write.writeString(item.Id); err != nil {
		return err
	}
	if err := write(item.Price); err != nil {
		return err
	}
	if err := write(int16(item.Quantity)); err != nil {
		return err
	}
	return write.writeString(item.Note)
}

func (InvMarshaler) UnmarshalInvoices(reader io.Reader) ([]*model.Invoice, error) {
	if err := checkInvVersion(reader); err != nil {
		return nil, err
	}
	count, err := readIntFromInt32(reader)
	if err != nil {
		return nil, err
	}
	invoices := make([]*model.Invoice, 0, count)
	for i := 0; i < count; i++ {
		invoice, err := readInvInvoice(reader)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func checkInvVersion(reader io.Reader) error {
	var magic uint32
	if err := binary.Read(reader, byteOrder, &magic); err != nil {
		return err
	}
	if magic != constants.MagicNumber {
		return fmt.Errorf("cannot read non-invoices inv file")
	}
	var version uint16
	if err := binary.Read(reader, byteOrder, &version); err != nil {
		return err
	}
	if version > constants.FileVersion {
		return fmt.Errorf("version %d is too new to read", version)
	}
	return nil
}
func readIntFromInt32(reader io.Reader) (int, error) {
	var i32 int32
	err := binary.Read(reader, byteOrder, &i32)
	return int(i32), err
}
func readInvInvoice(reader io.Reader) (invoice *model.Invoice, err error) {
	invoice = &model.Invoice{}
	for _, pId := range []*int{&invoice.Id, &invoice.CustomerId} {
		if *pId, err = readIntFromInt32(reader); err != nil {
			return nil, err
		}
	}
	for _, pDate := range []*time.Time{&invoice.Raised, &invoice.Due} {
		if *pDate, err = readInvDate(reader); err != nil {
			return nil, err
		}
	}
	if invoice.Paid, err = readBoolFromInt8(reader); err != nil {
		return nil, err
	}
	if invoice.Note, err = readInvString(reader); err != nil {
		return nil, err
	}
	var count int
	if count, err = readIntFromInt32(reader); err != nil {
		return nil, err
	}
	invoice.Items, err = readInvItems(reader, count)
	return invoice, err
}

func readInvDate(reader io.Reader) (time.Time, error) {
	var n int32
	if err := binary.Read(reader, byteOrder, &n); err != nil {
		return time.Time{}, err
	}
	return time.Parse(invDateFormat, fmt.Sprint(n))
}

func readBoolFromInt8(reader io.Reader) (bool, error) {
	var i8 int8
	err := binary.Read(reader, byteOrder, &i8)
	return i8 == 1, err
}

func readInvString(reader io.Reader) (string, error) {
	var lenght int32
	if err := binary.Read(reader, byteOrder, &lenght); err != nil {
		return "", err
	}
	raw := make([]byte, lenght)
	if err := binary.Read(reader, byteOrder, &raw); err != nil {
		return "", nil
	}
	return string(raw), nil
}

func readInvItems(reader io.Reader, count int) ([]*model.Item, error) {
	items := make([]*model.Item, 0, count)
	for i := 0; i < count; i++ {
		item, err := readInvItem(reader)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func readInvItem(reader io.Reader) (item *model.Item, err error) {
	item = &model.Item{}
	if item.Id, err = readInvString(reader); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, byteOrder, &item.Price); err != nil {
		return nil, err
	}
	if item.Quantity, err = readIntFromInt16(reader); err != nil {
		return nil, err
	}

	item.Note, err = readInvString(reader)

	return item, nil
}

func readIntFromInt16(reader io.Reader) (int, error) {
	var i16 int16
	err := binary.Read(reader, byteOrder, &i16)
	return int(i16), err
}

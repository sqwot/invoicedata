package model

import "encoding/xml"

type Item struct {
	Id       string
	Price    float64
	Quantity int
	Note     string
}

type XMLItem struct {
	XMLName  xml.Name `xml:"ITEM"`
	Id       string   `xml:",attr"`
	Price    float64  `xml:",attr"`
	Quantity int      `xml:",attr"`
	Note     string   `xml:"NOTE"`
}

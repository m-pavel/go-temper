package temper

import (
	"math"

	"github.com/google/gousb"
)

const (
	VENDOR_ID  gousb.ID = 0x1130
	PRODUCT_ID gousb.ID = 0x660c
)

const (
	CMD1 = "\x0a\x0b\x0c\x0d\x00\x00\x02\x00"
	CMD0 = "\x00\x00\x00\x00\x00\x00\x00\x00"
	CMD2 = "\x0a\x0b\x0c\x0d\x00\x00\x01\x00"
)

type Readings struct {
	Temp float64
	Rh   float64
}

func (r Readings) Dew() float64 {
	var tn, m float64
	if r.Temp > 0 {
		tn = 243.12
		m = 17.62
	} else {
		tn = 272.62
		m = 22.46
	}
	return tn * (math.Log(r.Rh/100) + (m * r.Temp / (tn + r.Temp))) / (m - math.Log(r.Rh/100) - (m * r.Temp / (tn + r.Temp)))
}

type Temper interface {
	Read() (*Readings, error)
	Close() error
}

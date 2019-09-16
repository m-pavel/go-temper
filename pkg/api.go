package temper

import "math"

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

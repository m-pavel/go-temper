package temper

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ltemper -lusb
// #include <temper.h>
// #include <stdlib.h>
// #include <usb.h>
import "C"
import (
	"log"
	"math"
)

type Temper struct {
	t *C.Temper
}

type Readings struct {
	Temp float64
	Rh   float64
}

func New(devicenum, timeout int, debug bool) (*Temper, error) {
	t := Temper{}
	var err error
	cdbg := 0
	if debug {
		cdbg = 1
	}

	_, err = C.usb_set_debug(C.int(cdbg))
	if err != nil {
		if debug {
			log.Println(err)
		}
	}
	_, err = C.usb_init()
	if err != nil {
		if debug {
			log.Println(err)
		}
	}
	_, err = C.usb_find_busses()
	if err != nil {
		if debug {
			log.Println(err)
		}
	}
	_, err = C.usb_find_devices()
	if err != nil {
		if debug {
			log.Println(err)
		}
	}

	t.t, err = C.TemperCreateFromDeviceNumber(C.int(devicenum), C.int(timeout*1000), C.int(cdbg))
	if err != nil {
		if debug {
			log.Println(err)
		}
	}
	if t.t == nil {
		return nil, err
	}
	return &t, nil
}

func (t *Temper) Close() error {
	_, err := C.TemperFree(t.t)
	return err
}

func (t *Temper) Read() (*Readings, error) {
	var tm, h C.double
	_, err := C.TemperGetTempAndRelHum(t.t, &tm, &h)
	if err != nil {
		return nil, err
	}
	return &Readings{Temp: float64(tm), Rh: float64(h)}, nil
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

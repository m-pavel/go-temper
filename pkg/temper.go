package temper
// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ltemper -lusb
// #include <temper.h>
// #include <stdlib.h>
// #include <usb.h>
import "C"

type Temper struct {
	t *C.Temper
}

type Readings struct {
	Temp float32
	Rh float32
}

func New(devicenum, timeout int ) (*Temper, error) {
	t := Temper{}
	var err error
	_, err = C.usb_set_debug(0);
	if err != nil {
		return nil, err
	}
	_, err = C.usb_init();
	if err != nil {
		return nil, err
	}
	_, err = C.usb_find_busses();
	if err != nil {
		return nil, err
	}
	_, err = C.usb_find_devices();
	if err != nil {
		return nil, err
	}

	t.t, err = C.TemperCreateFromDeviceNumber(C.int(devicenum), C.int(timeout * 1000), 0)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *Temper) Close() error{
	_, err := C.TemperFree(t.t)
	return err
}

func (t *Temper) Read() (*Readings, error) {
	var tm, h C.double
	_, err := C.TemperGetTempAndRelHum(t.t, &tm, &h)
	if (err != nil) {
		return nil, err
	}
	return &Readings{Temp: float32(tm), Rh: float32(h)}, nil
}

package temper
// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ltemper
// #include <temper.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"strconv"
)

type Temper struct {
	t *C.Temper
}

type Readings struct {
	Temp float32
	Rh float32
}

func New(devicenum, timeout int ) *Temper {
	t := Temper{}
	t.t = C.TemperCreateFromDeviceNumber(C.int(devicenum), C.int(timeout), 1)
	if t.t == nil {
		return nil
	}
	return &t
}

func (t *Temper) Close() {
	C.TemperFree(t.t)
}

func (t *Temper) Read() (*Readings, error) {
	var tm, h C.double
	res := C.TemperGetTempAndRelHum(t.t, &tm, &h)
	if (res != 0) {
		return nil, errors.New(strconv.Itoa(int(res)))
	}
	return &Readings{}, nil
}


/*
Temper *TemperCreateFromDeviceNumber(int deviceNum, int timeout, int debug);
void TemperFree(Temper *t);

int TemperGetTemperatureInC(Temper *t, double *tempC);
int TempterGetOtherStuff(Temper *t, char *buf, int length);

int TemperGetTempAndRelHum(Temper *t, double *tempC, double *relhum);

 */
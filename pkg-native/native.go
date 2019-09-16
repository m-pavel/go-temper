package tempern

import (
	"errors"

	"fmt"

	"time"

	"github.com/google/gousb"
	"github.com/m-pavel/go-temper/pkg"
)

type nTemper struct {
	dev   *gousb.Device
	ctx   *gousb.Context
	debug bool
}

const (
	VENDOR_ID  gousb.ID = 0x1130
	PRODUCT_ID gousb.ID = 0x660c
)

func New(devicenum, timeout int, debug bool) (temper.Temper, error) {
	nt := nTemper{debug: debug}
	nt.ctx = gousb.NewContext()

	if debug {
		nt.ctx.Debug(3)
	}

	var err error
	nt.dev, err = nt.ctx.OpenDeviceWithVIDPID(VENDOR_ID, PRODUCT_ID)
	if err != nil {
		return nil, err
	}
	if nt.dev == nil {
		return nil, errors.New("No TEMPERHum device found.")
	}
	return &nt, nil
}

func (t *nTemper) Close() error {
	if t.dev != nil {
		t.dev.Close()
		t.dev = nil
		return t.ctx.Close()
	}
	return nil
}

func (t *nTemper) Read() (*temper.Readings, error) {

	if se := t.sendCommand(10, 11, 12, 13, 0, 0, 2, 0); se != nil {
		return nil, se
	}
	if se := t.sendCommand(0x48, 0, 0, 0, 0, 0, 0, 0); se != nil {
		return nil, se
	}
	for i := 0; i < 7; i++ {
		if se := t.sendCommand(0, 0, 0, 0, 0, 0, 0, 0); se != nil {
			return nil, se
		}
	}
	if se := t.sendCommand(10, 11, 12, 13, 0, 0, 1, 0); se != nil {
		return nil, se
	}

	time.Sleep(400 * time.Millisecond)

	return t.getData()
}

func (t *nTemper) sendCommand(v ...byte) error {
	if t.debug {
		fmt.Printf("sending bytes %02x, %02x, %02x, %02x, %02x, %02x, %02x, %02x\n", v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7])
	}
	_, err := t.dev.Control(0x21, 9, 0x200, 0x01, v)
	return err
}

func (t *nTemper) getData() (*temper.Readings, error) {
	buf := make([]byte, 256)
	n, err := t.dev.Control(0xa1, 1, 0x300, 0x01, buf)
	if err != nil {
		return nil, err
	}
	if t.debug {
		fmt.Printf("Read %d bytes: %v\n", n, buf)
	}

	readings := temper.Readings{}
	temperature := (uint16(buf[1]) & 0xFF) + (uint16(buf[0]) << 8)
	fmt.Println(temperature)
	readings.Temp = -39.7 + .01*float64(temperature)

	rh := (buf[3] & 0xFF) + ((buf[2] & 0xFF) << 8)
	thum := -2.0468 + 0.0367*float64(rh) - 1.5955e-6*float64(rh)*float64(rh)
	readings.Rh = (readings.Temp-25)*(.01+.00008*float64(rh)) + thum
	return &readings, nil
}

package tempern

import (
	"errors"

	"fmt"

	"time"

	"log"

	"github.com/google/gousb"
	"github.com/m-pavel/go-temper/pkg"
)

type nTemper struct {
	dev *gousb.Device
	ctx *gousb.Context

	cfg      *gousb.Config
	if1, if2 *gousb.Interface
	debug    bool
}

func New(devicenum, timeout int, debug bool) (temper.Temper, error) {
	nt := nTemper{debug: debug}
	nt.ctx = gousb.NewContext()

	if debug {
		nt.ctx.Debug(3)
	}

	var err error
	if nt.dev, err = nt.ctx.OpenDeviceWithVIDPID(temper.VENDOR_ID, temper.PRODUCT_ID); err != nil {
		return nil, err
	}
	if nt.dev == nil {
		return nil, errors.New("No TEMPERHum device found.")
	}

	if err = nt.dev.SetAutoDetach(true); err != nil {
		nt.dev.Close()
		return nil, err
	}
	if nt.cfg, err = nt.dev.Config(1); err != nil {
		nt.Close()
		return nil, err
	}
	if nt.if1, err = nt.cfg.Interface(0, 0); err != nil {
		nt.Close()
		return nil, err
	}
	if nt.if2, err = nt.cfg.Interface(1, 0); err != nil {
		nt.Close()
		return nil, err
	}

	if _, err = nt.read(0x52); err != nil {
		nt.Close()
		return nil, err
	}
	return &nt, nil
}

func (t *nTemper) Close() error {
	if t.if1 != nil {
		t.if1.Close()
		t.if1 = nil
	}
	if t.if2 != nil {
		t.if2.Close()
		t.if2 = nil
	}
	if t.cfg != nil {
		if err := t.cfg.Close(); err != nil {
			log.Println(err)
		}
		t.cfg = nil
	}
	if t.dev != nil {
		if err := t.dev.Close(); err != nil {
			log.Println(err)
		}
		t.dev = nil
	}
	if t.ctx != nil {
		if err := t.ctx.Close(); err != nil {
			log.Println(err)
		}
		t.ctx = nil
	}
	return nil
}

func (t *nTemper) Read() (*temper.Readings, error) {
	return t.read(0x48)
}

func (t *nTemper) read(cb byte) (*temper.Readings, error) {

	if se := t.sendCommand([]byte(temper.CMD1)); se != nil {
		return nil, se
	}
	req := []byte(temper.CMD0)
	req[0] = cb
	if se := t.sendCommand(req); se != nil {
		return nil, se
	}
	for i := 0; i < 7; i++ {
		if se := t.sendCommand([]byte(temper.CMD0)); se != nil {
			return nil, se
		}
	}
	if se := t.sendCommand([]byte(temper.CMD2)); se != nil {
		return nil, se
	}

	time.Sleep(600 * time.Millisecond)

	return t.getData()
}

func (t *nTemper) sendCommand(v []byte) error {
	if t.debug {
		fmt.Printf("sending bytes %02x, %02x, %02x, %02x, %02x, %02x, %02x, %02x\n", v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7])
	}
	_, err := t.dev.Control(0x21, 9, 0x200, 0x01, v)
	return err
}

func (t *nTemper) getData() (*temper.Readings, error) {
	buf := make([]byte, 256)
	if n, err := t.dev.Control(0xa1, 1, 0x300, 0x01, buf); err != nil {
		return nil, err
	} else if t.debug {
		fmt.Printf("Read %d bytes: %v\n", n, buf)
	}

	readings := temper.Readings{}
	temperature := (uint16(buf[1]) & 0xFF) + (uint16(buf[0]) << 8)
	readings.Temp = -39.7 + .01*float64(temperature)

	rh := (uint16(buf[3]) & 0xFF) + ((uint16(buf[2]) & 0xFF) << 8)
	thum := -2.0468 + 0.0367*float64(rh) - 1.5955e-6*float64(rh)*float64(rh)
	readings.Rh = (readings.Temp-25)*(.01+.00008*float64(rh)) + thum
	return &readings, nil
}

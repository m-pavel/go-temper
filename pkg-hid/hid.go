package temperhid

import (
	"fmt"

	"github.com/GeertJohan/go.hid"
	"github.com/m-pavel/go-temper/pkg"
)

type hTemper struct {
	hd *hid.Device
}

func New(devicenum, timeout int, debug bool) (temper.Temper, error) {
	ht := hTemper{}
	var err error
	ht.hd, err = hid.Open(uint16(temper.VENDOR_ID), uint16(temper.PRODUCT_ID), "")
	return &ht, err
}

func (t *hTemper) Close() error {
	t.hd.Close()
	return nil
}

func (t *hTemper) Read() (*temper.Readings, error) {
	if _, err := t.hd.Write([]byte(temper.CMD1)); err != nil {
		return nil, err
	}

	req := []byte(temper.CMD0)
	req[0] = 0x48
	if _, err := t.hd.Write([]byte(req)); err != nil {
		return nil, err
	}

	for i := 0; i < 7; i++ {
		if _, err := t.hd.Write([]byte(temper.CMD0)); err != nil {
			return nil, err
		}
	}
	if _, err := t.hd.Write([]byte(temper.CMD2)); err != nil {
		return nil, err
	}

	buffer := make([]byte, 8)
	n, err := t.hd.Read(buffer)
	if err != nil {
		return nil, err
	}

	fmt.Println(n)
	fmt.Println(buffer)
	return nil, nil
}

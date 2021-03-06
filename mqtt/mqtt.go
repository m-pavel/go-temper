package main

import (
	"github.com/m-pavel/go-hassio-mqtt/pkg"
	"github.com/m-pavel/go-temper/pkg"
	"github.com/m-pavel/go-temper/pkg-native"
)

type TemperMqtt struct {
	Temp float64 `json:"temperature"`
	Rh   float64 `json:"humidity"`
}

type TemperService struct {
	ghm.NonListerningService
	t temper.Temper
}

func (ts TemperService) PrepareCommandLineParams() {}
func (ts TemperService) Name() string              { return "temper" }

func (ts *TemperService) Init(ctx *ghm.ServiceContext) error {
	var err error
	ts.t, err = tempern.New(0, 0, ctx.Debug())
	return err
}

func (ts TemperService) Do() (interface{}, error) {
	rd, err := ts.t.Read()
	if err != nil {
		return nil, err
	}

	return &TemperMqtt{Temp: rd.Temp, Rh: rd.Rh}, nil
}

func (ts TemperService) Close() error {
	return ts.t.Close()
}

func main() {
	ghm.NewStub(&TemperService{}).Main()
}

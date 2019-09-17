package main

import (
	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/m-pavel/go-hassio-mqtt/pkg"
	"github.com/m-pavel/go-temper/pkg"
	"github.com/m-pavel/go-temper/pkg-native"
)

type TemperMqtt struct {
	Temp float64 `json:"temperature"`
	Rh   float64 `json:"humidity"`
}

type TemperService struct {
	t     temper.Temper
	topic string
}

func (ts TemperService) PrepareCommandLineParams() {}
func (ts TemperService) Name() string              { return "temper" }

func (ts *TemperService) Init(client MQTT.Client, topic, topicc, topica string, debug bool) error {
	var err error
	ts.t, err = tempern.New(0, 0, debug)
	ts.topic = topic
	return err
}

func (ts TemperService) Do(client MQTT.Client) error {
	rd, err := ts.t.Read()
	if err != nil {
		return err
	}

	mqt := TemperMqtt{Temp: rd.Temp, Rh: rd.Rh}
	bp, err := json.Marshal(&mqt)
	if err != nil {
		return err
	}
	tkn := client.Publish(ts.topic, 0, false, bp)
	return tkn.Error()
}

func (ts TemperService) Close() error {
	return ts.t.Close()
}

func main() {
	hs := TemperService{}
	hmss := ghm.NewStub(&hs)
	hmss.Main()
}

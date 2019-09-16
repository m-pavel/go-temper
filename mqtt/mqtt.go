package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/m-pavel/go-temper/pkg-c"
	"github.com/sevlyar/go-daemon"
)

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func main() {
	var logf = flag.String("log", "influxtemper.log", "log")
	var pid = flag.String("pid", "influxtemper.pid", "pid")
	var notdaemonize = flag.Bool("n", false, "Do not do to background.")
	var signal = flag.String("s", "", `send signal to the daemon stop — shutdown`)
	var iserver = flag.String("mqtt", "tcp://localhost:1883", "MQTT endpoint")
	var idb = flag.String("t", "nn/temper", "MQTT topic")
	var debug = flag.Bool("d", false, "debug")
	var interval = flag.Int("interval", 10, "Interval secons")
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)

	cntxt := &daemon.Context{
		PidFileName: *pid,
		PidFilePerm: 0644,
		LogFileName: *logf,
		LogFilePerm: 0640,
		WorkDir:     "/tmp",
		Umask:       027,
		Args:        os.Args,
	}

	if !*notdaemonize && len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %v", err)
		}
		daemon.SendCommands(d)
		return
	}

	if !*notdaemonize {
		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal(err)
		}
		if d != nil {
			return
		}
	}

	daemonf(*iserver, *idb, *interval, *debug)
}

type TemperMqtt struct {
	Temp float64
	Rh   float64
}

func daemonf(iserver, db string, interval int, debug bool) {
	var err error
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: iserver,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	t, err := temperc.New(0, 5, debug)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	failcnt := 0

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("temper-go-cli")
	topic := "nn/sensors"
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		if failcnt >= 15 {
			return
		}
		rd, err := t.Read()
		mqt := TemperMqtt{Temp: rd.Temp, Rh: rd.Rh}
		if err != nil {
			log.Println(err)
			failcnt += 1
		} else {
			failcnt = 0
			client.Publish(topic, 0, false, &mqt)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

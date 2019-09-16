package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/m-pavel/go-temper/pkg-native"
	"github.com/sevlyar/go-daemon"
)

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func main() {
	var logf = flag.String("log", "tempermqtt.log", "log")
	var pid = flag.String("pid", "tempermqtt.pid", "pid")
	var notdaemonize = flag.Bool("n", false, "Do not do to background.")
	var signal = flag.String("s", "", `send signal to the daemon stop — shutdown`)
	var mqtt = flag.String("mqtt", "tcp://localhost:1883", "MQTT endpoint")
	var topic = flag.String("t", "nn/temper", "MQTT topic")
	var user = flag.String("mqtt-user", "", "MQTT user")
	var pass = flag.String("mqtt-pass", "", "MQTT password")
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

	daemonf(*mqtt, *topic, *user, *pass, *interval, *debug)
}

type TemperMqtt struct {
	Temp float64
	Rh   float64
}

func daemonf(mqtt, topic string, u, p string, interval int, debug bool) {
	var err error

	t, err := tempern.New(0, 5, debug)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	failcnt := 0

	opts := MQTT.NewClientOptions().AddBroker(mqtt)
	opts.SetClientID("temper-go-cli")
	if u != "" {
		opts.Username = u
		opts.Password = p
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		if failcnt >= 15 {
			return
		}
		rd, err := t.Read()
		if err != nil {
			log.Println(err)
			failcnt += 1
		} else {
			if debug {
				log.Println(rd)
			}
			mqt := TemperMqtt{Temp: rd.Temp, Rh: rd.Rh}

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

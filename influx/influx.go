package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/m-pavel/go-temper/pkg"
	"github.com/sevlyar/go-daemon"
)

func main() {
	var logf = flag.String("log", "influxtemper.log", "log")
	var pid = flag.String("pid", "influxtemper.pid", "pid")
	var notdaemonize = flag.Bool("n", false, "Do not do to background.")
	var signal = flag.String("s", "", `send signal to the daemon stop — shutdown`)
	var iserver = flag.String("influx", "http://localhost:8086", "Influx DB endpoint")
	var idb = flag.String("db", "tion", "Influx DB")

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

	daemonf(*iserver, *idb, *interval)
}

func daemonf(iserver, db string, interval int) {
	var err error
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: iserver,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	t, err := temper.New(0, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	for {
		rd, err := t.Read()
		if err != nil {
			log.Println(err)
		} else {
			point, err := client.NewPoint(db,
				map[string]string{},
				map[string]interface{}{
					"temper_t": rd.Temp,
					"temper_h": rd.Rh,
					"temper_d": 0,
				},
				time.Now())
			if err != nil {
				log.Printf("Insert data error: %v", err)
				return
			}
			bp, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "tion",
				Precision: "s",
			})
			if err != nil {
				log.Printf("Insert data error: %v", err)
				return
			}
			bp.AddPoint(point)
			err = cli.Write(bp)
			if err != nil {
				log.Printf("Insert data error: %v", err)
				return
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

	done <- struct{}{}
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

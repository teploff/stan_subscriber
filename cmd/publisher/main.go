package main

import (
	"encoding/json"
	"flag"
	"github.com/nats-io/stan.go"
	"github.com/teploff/stan_subscriber/config"
	"github.com/teploff/stan_subscriber/domain"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFile = flag.String("config", "./init/config_dev.yaml", "configuration file path")
)

func main() {
	flag.Parse()

	cfg, err := config.LoadFromFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	connect, err := stan.Connect(cfg.Stan.ClusterID, cfg.Stan.ClientID,
		stan.Pings(60, 2*60),
		stan.NatsURL(cfg.Stan.Addr),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				msg, err := json.Marshal(domain.Measurement{
					Ts:        time.Now(),
					ActorName: "TestActor",
					Type:      "Some type",
					Data:      "Some data",
				})
				if err != nil {
					log.Fatalln(err)
				}
				_, err = connect.PublishAsync("measurements", msg, func(_ string, _ error) {})
				if err != nil {
					log.Fatalln(err)
				}
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	if err = connect.Close(); err != nil {
		log.Fatal(err)
	}
}

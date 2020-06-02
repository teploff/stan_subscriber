package main

import (
	"flag"
	"github.com/nats-io/stan.go"
	"github.com/teploff/stan_subscriber/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile = flag.String("config", "./init/config.yaml", "configuration file path")
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

	var sub stan.Subscription
	go func() {
		sub, err = connect.QueueSubscribe(cfg.Stan.Subject,
			"q_group_1",
			func(m *stan.Msg) {
				if err = m.Ack(); err != nil {
					panic(err)
				}
				log.Printf("Сообщение =  %s\n", string(m.Data))
			},
			stan.SetManualAckMode(),
		)
		sub.SetPendingLimits(-1, -1)
		if err != nil {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	if err = sub.Close(); err != nil {
		log.Fatal(err)
	}
}

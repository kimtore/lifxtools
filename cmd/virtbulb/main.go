// Emulate a LIFX bulb and draw out its contents on screen.

package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/virtbulb"
	log "github.com/sirupsen/logrus"
)

const (
	EnvBindAddress = "BIND_ADDRESS"
)

func main() {
	err := run()
	if err != nil {
		log.Errorf("fatal: %s", err)
	}
}

func run() error {
	bindAddress := flag.String("bind", getEnvDefault(EnvBindAddress, "127.0.0.1:56700"), "Listen address for the virtual LIFX bulb")

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})

	log.SetLevel(log.DebugLevel)

	log.Infof("Virtual Bulb starting up")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	udp, err := net.ListenPacket("udp", *bindAddress)
	if err != nil {
		return err
	}

	stream := virtbulb.NewPacketStreamer(ctx, udp)

	cv := virtbulb.NewCanvas()
	drawTimer := time.NewTicker(time.Second)

	log.Infof("Setup complete")

	for {
		select {
		case <-ctx.Done():
			log.Infof("Virtual Bulb shutting down.")
			return nil
		case s := <-sigs:
			log.Infof("Received %s signal", s.String())
			cancel()
		case <-drawTimer.C:
			cv.Print()
		case p := <-stream.C:
			err := virtbulb.PaintPacket(cv, p)
			if err != nil {
				log.Debug(err)
			}
		}
	}
}

func getEnvDefault(env, fallback string) string {
	value, found := os.LookupEnv(env)
	if found {
		return value
	}
	return fallback
}

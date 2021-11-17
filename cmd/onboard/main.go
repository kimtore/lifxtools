package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := run()
	if err != nil {
		log.Errorf("fatal: %s", err)
		os.Exit(1)
	}
}

var (
	addr = flag.String("addr", "127.16.0.1:56700", "Network address of LIFX bulb")
	ssid = flag.String("ssid", "", "Wifi SSID")
	psk  = flag.String("psk", "", "Wifi pre-shared key")
)

func run() error {
	flag.Parse()

	if len(*ssid) == 0 || len(*psk) == 0 {
		return fmt.Errorf("missing parameters")
	}

	payload := &lifx.SetAccessPointMessage{
		Reserved1: 0x02,
		SSID:      *ssid,
		PSK:       *psk,
		Security:  lifx.Security_WPA2_AES_PSK,
	}
	packet := lifx.NewPacket(payload)

	buf := &bytes.Buffer{}
	err := packet.Write(buf)
	if err != nil {
		return err
	}

	log.Infof("Connecting to TCP %s...", *addr)
	sock, err := tls.DialWithDialer(
		&net.Dialer{
			Timeout: time.Second * 5,
		},
		"tcp", *addr,
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		return err
	}
	defer sock.Close()

	log.Infof("Sending onboarding payload...")
	_, err = sock.Write(buf.Bytes())
	if err != nil {
		return err
	}

	log.Infof("Done, happy hacking.")

	return nil
}

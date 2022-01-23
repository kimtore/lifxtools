package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/dorkowscy/lyslix/lifx"
)

func main() {
	host := flag.String("host", "", "lifx bulb hostname or ip address")
	port := flag.Uint("port", 56700, "lifx bulb port")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client, err := lifx.NewClient(addr)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	framerate := time.Second / 10

	rand.Seed(time.Now().UnixNano())

	for hue := uint16(0); hue < lifx.MaxHue; hue += 50 {
		color := lifx.HBSK{
			Hue:        hue,
			Saturation: lifx.FullSaturation,
			Brightness: lifx.FullBrightness,
		}
		fmt.Println(color)
		client.SetColor(color, 0)
		time.Sleep(framerate)
	}
}

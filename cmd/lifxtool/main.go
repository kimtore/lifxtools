package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	host := flag.String("host", "", "lifx bulb hostname or ip address")
	port := flag.Uint("port", 56700, "lifx bulb port")
	size := flag.Uint("size", 0, "lifx canvas size")
	fps := flag.Uint("fps", 3, "canvas updates per second")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client, err := lifx.NewClient(addr)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	framerate := time.Second / time.Duration(*fps)
	cv := canvas.New(client, int(*size))
	cv.Clear()

	pixels := cv.Pixels()
	eff := &effects.NorthernLights{
		Threshold: 0.05,
		Cutoff:    0.5,
		Base:      colorful.Hsl(0, 0.8, 0),
		Color1:    colorful.Hsl(60, 1, 1),
		Color2:    colorful.Hsl(240, 1, 1),
		Fade:      0.5,
	}

	for {
		eff.Draw(pixels)
		cv.Draw(framerate)
		time.Sleep(framerate)
	}
}

package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

const rad = math.Pi / 180

func main() {
	host := flag.String("host", "", "lifx bulb hostname or ip address")
	port := flag.Uint("port", 56700, "lifx bulb port")
	size := flag.Uint("size", 0, "lifx canvas size")
	fps := flag.Uint("fps", 3, "canvas updates per second")
	flag.Parse()

	framerate := time.Second / time.Duration(*fps)
	_ = size

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client, err := lifx.NewClient(addr)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	cv := canvas.NewStrip(client, int(*size))
	base := colorful.Hcl(0, 0.25, 0.25)
	base=colorful.Color{}
	cv.Fill(base)
	cv.Draw(0)

	pixels := cv.Pixels()

	eff := &effects.NorthernLights{
		Threshold: 0.02,
		Cutoff:    0.3,
		Base:      base,
		Color1:    colorful.Hcl(220, 0.25, 0.75),
		Color2:    colorful.Hcl(300, 0.25, 0.75),
		Fade:      0.2,
	}

	for {
		eff.Draw(pixels)
		cv.Set(pixels)
		cv.Draw(framerate)
		time.Sleep(framerate)
		fmt.Println()
	}
}

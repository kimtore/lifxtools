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

func circle(cv canvas.Canvas, chroma, lum float64, interval time.Duration) {
	fmt.Printf("L*,C*,h*,H,S,V\n")
	for hue := 0.0; hue < 360; hue++ {
		/*
			rads := deg * rad
			x := math.Sin(rads)
			y := math.Cos(rads)
			colorful.Hcl()
			color := colorful.Lab(lum, x, y)

		*/
		color := colorful.Hcl(hue, chroma, lum)
		clamped := color.Clamped()
		h, s, v := color.Hsv()
		fmt.Printf("%.5f,%.5f,%.5f,%.5f,%.5f,%.5f", lum, chroma, hue, h, s, v)
		if clamped != color {
			fmt.Printf(" [out of gamut]")
		}
		fmt.Printf("\n")
		//fmt.Printf("%.5f,%.5f,%.5f,%.5f,%.5f,%.5f,%.5f\n", deg, rads, x, y, h, s, v)
		//log.Printf("(%.5f, %.5f) = (%.5f, %.5f, %.5f)", x, y, h, s, v)
		cv.Fill(color)
		cv.Draw(interval)
		time.Sleep(interval)
	}
}

package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
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

	//framerate := time.Second / 10
	cv := canvas.New(client, 24)
	cv.Clear()

	pixels := cv.Pixels()
	for i := range pixels {
		pixels[i] = colorful.Hsl(rand.Float64()*360, rand.Float64(), rand.Float64())
	}

	fmt.Println(pixels)

	cv.Draw()
}

// Emulate a LIFX bulb and draw out its contents on screen.

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/dorkowscy/lifxtool/pkg/virtbulb"
	"github.com/gdamore/tcell"
	"github.com/lucasb-eyer/go-colorful"
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
	log.SetOutput(os.Stderr)

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
	drawTimer := time.NewTicker(time.Millisecond * 10)

	statsUpdate := time.Second
	statsTimer := time.NewTicker(statsUpdate)

	screen, err := initTcell()
	if err != nil {
		return err
	}
	defer screen.Fini()

	// Listen for keyboard input
	consoleEvents := make(chan tcell.Event, 1)
	go func() {
		for {
			consoleEvents <- screen.PollEvent()
		}
	}()

	log.Infof("Setup complete")

	screen.Clear()
	packetsReceived := 0
	lastPackets := 0

	for {
		select {
		case <-ctx.Done():
			log.Infof("Virtual Bulb shutting down.")
			return nil
		case s := <-sigs:
			log.Infof("Received %s signal", s.String())
			cancel()
		case <-statsTimer.C:
			diff := float64(packetsReceived - lastPackets)
			lastPackets = packetsReceived
			drawInfo(screen, 0, diff*float64(statsUpdate)/float64(time.Second))
		case <-drawTimer.C:
			drawCanvas(screen, cv.Pixels())
		case p := <-stream.C:
			packetsReceived++
			err := virtbulb.PaintPacket(cv, p)
			if err != nil {
				log.Debug(err)
			}
		case ev := <-consoleEvents:
			switch ev := ev.(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					cancel()
				}
			}
		}
	}
}

func drawInfo(s tcell.Screen, y int, packetsPerSecond float64) {
	const packetHigh = 30.0
	const packetMedium = 20.0
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	switch {
	case packetsPerSecond > packetHigh:
		style = style.Background(tcell.ColorRed)
	case packetsPerSecond > packetMedium:
		style = style.Background(tcell.ColorOrange)
	default:
		style = style.Background(tcell.ColorDarkGreen)
	}
	drawText(s, 1, y, 50, y, style,
		fmt.Sprintf("%-7.1f packets per second", packetsPerSecond))
}

func drawCanvas(s tcell.Screen, pixels []colorful.Color) {
	for i := range pixels {
		drawLine(s, pixels[i], i+2)
	}
	s.Show()
}

func drawLine(s tcell.Screen, c colorful.Color, line int) {
	const asciiBlock = rune(219)
	r, g, b := c.Clamped().RGB255()
	color := tcell.NewRGBColor(int32(r), int32(g), int32(b))
	style := tcell.StyleDefault.Background(color).Foreground(color)
	s.SetContent(1, line, asciiBlock, nil, style)
	s.SetContent(2, line, asciiBlock, nil, style)
	hcl := textutil.SprintfHCL(c)
	drawText(s, 5, line, 40, line, tcell.StyleDefault, hcl)
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func initTcell() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = s.Init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func getEnvDefault(env, fallback string) string {
	value, found := os.LookupEnv(env)
	if found {
		return value
	}
	return fallback
}

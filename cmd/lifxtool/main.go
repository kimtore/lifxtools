package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/dorkowscy/lifxtool/pkg/server"
	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	err := run()
	if err != nil {
		log.Errorf("fatal: %s", err)
	}
}

func run() error {
	configFileName := flag.String("config", "config.yaml", "path to configuration file")
	presetName := flag.String("preset", "", "preset effect and canvas combination to run")
	bindAddress := flag.String("bind", "127.0.0.1:7178", "HOST:PORT combination for setting up an HTTP server")

	flag.Parse()

	cfg, err := readConfig(*configFileName)
	if err != nil {
		return err
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})

	if cfg.Options.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("LIFXTOOL starting up")

	log.Infof("Initializing bulbs...")
	bulbs, err := bulbClients(cfg.Bulbs)
	if err != nil {
		return err
	}

	log.Infof("Initializing canvases...")
	canvases, err := buildCanvases(cfg.Canvases, bulbs)
	if err != nil {
		return err
	}

	presets := mapPresets(cfg.Presets)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mgr := effects.NewManager(ctx, presets, canvases, bulbs)
	srv := server.New(mgr)

	// Start one-off with preset if requested on command line.
	if len(*presetName) > 0 {
		err = mgr.StartPreset(*presetName)
		if err != nil {
			return err
		}
	}

	// Start a HTTP server that controls the effects.
	if len(*bindAddress) > 0 {
		go func() {
			log.Infof("Starting HTTP server on %s...", *bindAddress)
			err := http.ListenAndServe(*bindAddress, srv.Router())
			if err != nil {
				log.Errorf("HTTP server shut down with error: %s", err)
			}
			cancel()
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-sigs:
			return nil
		}
	}
}

func mapPresets(presets []config.Preset) map[string]config.Preset {
	mp := make(map[string]config.Preset)
	for _, preset := range presets {
		mp[preset.Name] = preset
	}
	return mp
}

func buildCanvases(cvs []config.Canvas, bulbs map[string]lifx.Client) (map[string]canvas.Canvas, error) {
	var err error
	canvases := make(map[string]canvas.Canvas)
	for i, cv := range cvs {
		log.Debugf("Initializing new canvas '%s'", cv.Name)
		canvases[cv.Name], err = buildCanvas(cv, bulbs)
		if err != nil {
			return nil, fmt.Errorf("canvas %d: %w", i, err)
		}
	}
	return canvases, nil
}

func buildCanvas(cv config.Canvas, bulbs map[string]lifx.Client) (canvas.Canvas, error) {
	seq := make([]canvas.Canvas, 0, len(bulbs))

	for _, bulb := range cv.Bulbs {
		client := bulbs[bulb.Name]
		if client == nil {
			return nil, fmt.Errorf("no bulb configured with name '%s'", bulb.Name)
		}
		if bulb.Zone.Min != nil && bulb.Zone.Max != nil {
			seq = append(seq, canvas.NewStrip(client, *bulb.Zone.Min-1, *bulb.Zone.Max-1))
			log.Debugf("Canvas '%s': using zones %d-%d from strip '%s'", cv.Name, *bulb.Zone.Min, *bulb.Zone.Max, bulb.Name)
		} else {
			seq = append(seq, canvas.NewLight(client))
			log.Debugf("Canvas '%s': using single color on bulb/strip '%s'", cv.Name, bulb.Name)
		}
	}

	return canvas.NewAggregate(seq...), nil
}

func bulbClients(bulbs []config.Bulb) (map[string]lifx.Client, error) {
	var err error
	clients := make(map[string]lifx.Client)
	for i, bulb := range bulbs {
		log.Debugf("Initializing new bulb with name='%s' and host='%s'", bulb.Name, bulb.Host)
		clients[bulb.Name], err = bulbClient(bulb)
		if err != nil {
			return nil, fmt.Errorf("bulb %d: %w", i, err)
		}
	}
	return clients, nil
}

func bulbClient(bulb config.Bulb) (lifx.Client, error) {
	const port = ":56700"
	if len(bulb.Name) == 0 {
		return nil, fmt.Errorf("missing bulb name")
	}
	if len(bulb.Host) == 0 {
		return nil, fmt.Errorf("missing hostname")
	}
	addr := bulb.Host + port
	return lifx.NewClient(addr)
}

/*
	base := colorful.Hcl(0, 0.25, 0.25)
	base = colorful.Color{}
	cv.Fill(base)
	cv.Draw(0)

	pixels := cv.Pixels()

	_ = &effects.NorthernLights{
	}

*/

func readConfig(filename string) (*config.Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	cfg := &config.Config{}
	r := yaml.NewDecoder(configFile)
	err = r.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	var eff effects.Effect

	configFileName := flag.String("config", "config.yaml", "path to configuration file")
	presetName := flag.String("preset", "", "preset effect and canvas combination to run")
	flag.Parse()

	cfg, err := readConfig(*configFileName)
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})

	if cfg.Options.Debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Initializing bulbs...")
	bulbs, err := bulbClients(cfg.Bulbs)
	if err != nil {
		panic(err)
	}

	log.Infof("Initializing canvases...")
	canvases, err := buildCanvases(cfg.Canvases, bulbs)
	if err != nil {
		panic(err)
	}

	log.Infof("Initializing presets...")
	preset, err := lookupPreset(cfg.Presets, *presetName)
	if err != nil {
		panic(err)
	}

	log.Infof("Using preset '%s'", preset.Name)
	cv := canvases[preset.Canvas]
	if cv == nil {
		err = fmt.Errorf("preset '%s' uses undefined canvas '%s'", preset.Name, preset.Canvas)
		panic(err)
	}

	log.Infof("Initializing effect '%s' on canvas '%s'", preset.Effect.Name, preset.Canvas)
	eff, err = makeEffect(preset.Effect)
	if err != nil {
		panic(err)
	}

	pixels := make([]colorful.Color, cv.Size())

	log.Debugf("Canvas contains a total of %d separate light points", len(pixels))

	for {
		log.Debugf("Tick()")
		eff.Draw(pixels)
		cv.Set(pixels)
		cv.Draw(preset.Delay)
		time.Sleep(preset.Delay)
	}

}

func lookupPreset(presets []config.Preset, name string) (*config.Preset, error) {
	for _, preset := range presets {
		if preset.Name == name {
			return &preset, nil
		}
	}
	return nil, fmt.Errorf("preset with name '%s' not found", name)
}

func makeEffect(cfg config.Effect) (effects.Effect, error) {
	var eff effects.Effect
	switch cfg.Name {
	case "colorwheel":
		eff = &effects.ColorWheel{}
	case "northernlights":
		eff = &effects.NorthernLights{}
	default:
		return nil, fmt.Errorf("effect '%s' does not exist", cfg.Name)
	}
	cf, err := json.Marshal(cfg.Config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cf, eff)
	if err != nil {
		return nil, err
	}
	return eff, nil
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

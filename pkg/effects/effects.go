package effects

import (
	"encoding/json"
	"fmt"

	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/lucasb-eyer/go-colorful"
)

type Effect interface {
	Init(pixels []colorful.Color)
	Draw(pixels []colorful.Color)
}

var library = map[string]Effect{}

func register(name string, effect Effect) {
	library[name] = effect
}

func fill(pixels []colorful.Color, color colorful.Color) {
	for i := range pixels {
		pixels[i] = color
	}
}

func NewFromConfig(cfg config.Effect) (Effect, error) {
	var eff Effect
	switch cfg.Name {
	case "colorwheel":
		eff = &ColorWheel{}
	case "northernlights":
		eff = &NorthernLights{}
	case "sine":
		eff = &Sine{}
	default:
		return nil, fmt.Errorf("effect '%s' is not implemented", cfg.Name)
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

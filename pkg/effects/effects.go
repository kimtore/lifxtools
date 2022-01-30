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

var library = map[string]func() Effect{}

func register(name string, factory func() Effect) {
	library[name] = factory
}

func fill(pixels []colorful.Color, color colorful.Color) {
	for i := range pixels {
		pixels[i] = color
	}
}

func NewFromConfig(cfg config.Effect) (Effect, error) {
	var eff Effect

	factory := library[cfg.Name]
	if factory == nil {
		return nil, fmt.Errorf("effect '%s' is not implemented", cfg.Name)
	}

	eff = factory()

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

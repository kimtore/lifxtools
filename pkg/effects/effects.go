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

func NewFromConfig(cfg config.Effect) (Effect, error) {
	var eff Effect
	switch cfg.Name {
	case "colorwheel":
		eff = &ColorWheel{}
	case "northernlights":
		eff = &NorthernLights{}
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

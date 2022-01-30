package runner

import (
	"context"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
)

type Manager interface {
	StartPreset(preset string) error
	StopPreset(preset string) error
	IsActive(preset string) bool
	Configuration(preset string) map[string]interface{}
	Configure(preset string, values map[string]interface{}) error
}

type manager struct {
	ctx      context.Context
	presets  map[string]config.Preset
	canvases map[string]canvas.Canvas
	bulbs    map[string]lifx.Client
	runners  map[string]*runner
}

func NewManager(
	ctx context.Context,
	presets map[string]config.Preset,
	canvases map[string]canvas.Canvas,
	bulbs map[string]lifx.Client,
) Manager {

	runners := make(map[string]*runner)
	for preset := range presets {
		eff, err := effects.NewFromConfig(presets[preset].Effect)
		if err != nil {
			log.Errorf("Initialize preset '%s': %s", preset, err)
			continue
		}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		run := &runner{
			ctx:    ctx,
			cancel: cancel,
			canvas: canvases[presets[preset].Canvas],
			delay:  presets[preset].Delay,
			effect: eff,
			name:   presets[preset].Name,
		}
		runners[preset] = run
	}

	return &manager{
		ctx:      ctx,
		presets:  presets,
		canvases: canvases,
		bulbs:    bulbs,
		runners:  runners,
	}
}

var _ Manager = &manager{}

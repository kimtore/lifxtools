package effects

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

type Effect interface {
	Draw(pixels []colorful.Color)
}

type Manager interface {
	StartPreset(preset string) error
	StopPreset(preset string) error
}

type manager struct {
	ctx      context.Context
	presets  map[string]config.Preset
	canvases map[string]canvas.Canvas
	bulbs    map[string]lifx.Client
	effects  map[string]Effect
	runners  map[string]*runner
}

func NewManager(
	ctx context.Context,
	presets map[string]config.Preset,
	canvases map[string]canvas.Canvas,
	bulbs map[string]lifx.Client,
	effects map[string]Effect,
) Manager {
	return &manager{
		ctx:      ctx,
		presets:  presets,
		canvases: canvases,
		bulbs:    bulbs,
		effects:  effects,
		runners:  make(map[string]*runner),
	}
}

var _ Manager = &manager{}

var (
	ErrPresetNotFound   = errors.New("preset not found")
	ErrPresetNotRunning = errors.New("preset not running")
)

func (r *manager) StartPreset(preset string) error {
	p, ok := r.presets[preset]
	if !ok {
		return ErrPresetNotFound
	}

	// Stop effect if it's already running
	_ = r.StopPreset(preset)

	eff, err := makeEffect(p.Effect)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(r.ctx)
	run := &runner{
		canvas: r.canvases[p.Canvas],
		cancel: cancel,
		ctx:    ctx,
		delay:  p.Delay,
		effect: eff,
		name:   p.Name,
	}

	r.runners[p.Name] = run
	go run.Run()

	return nil
}

func (r *manager) StopPreset(preset string) error {
	run := r.runners[preset]
	if run == nil {
		return ErrPresetNotRunning
	}

	run.cancel()
	delete(r.runners, preset)

	return nil
}

func makeEffect(cfg config.Effect) (Effect, error) {
	var eff Effect
	switch cfg.Name {
	case "colorwheel":
		eff = &ColorWheel{}
	case "northernlights":
		eff = &NorthernLights{}
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

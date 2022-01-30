package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
)

type Manager interface {
	StartPreset(preset string) error
	StopPreset(preset string) error
	Exists(preset string) bool
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

var (
	ErrPresetNotFound = errors.New("preset not found")
)

func (r *manager) IsActive(preset string) bool {
	run := r.runners[preset]
	if run == nil {
		return false
	}
	return run.ctx.Err() == nil
}

func (r *manager) StartPreset(preset string) error {
	run, ok := r.runners[preset]
	if !ok {
		return ErrPresetNotFound
	}

	// Stop effect if it's already running
	_ = r.StopPreset(preset)

	run.ctx, run.cancel = context.WithCancel(r.ctx)

	go run.Run()

	return nil
}

func (r *manager) StopPreset(preset string) error {
	run := r.runners[preset]
	if run == nil {
		return ErrPresetNotFound
	}
	run.cancel()
	return nil
}

func (r *manager) Exists(preset string) bool {
	return r.runners[preset] != nil
}

// Return the active runner's configuration for this preset if available,
// otherwise return the default configuration for that preset.
func (r *manager) Configuration(preset string) map[string]interface{} {
	run := r.runners[preset]
	if run == nil {
		return nil
	}
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(run.effect)
	values := make(map[string]interface{})
	json.NewDecoder(buf).Decode(&values)
	return values
}

func (r *manager) Configure(preset string, values map[string]interface{}) error {
	run := r.runners[preset]
	if run == nil {
		return nil
	}
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(values)
	return json.NewDecoder(buf).Decode(run.effect)
}

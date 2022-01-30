package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

type runner struct {
	canvas canvas.Canvas
	cancel context.CancelFunc
	ctx    context.Context
	delay  time.Duration
	effect effects.Effect
	name   string
}

func (r *runner) Run() {
	log.Infof("[%s] Runner started", r.name)

	pixels := make([]colorful.Color, r.canvas.Size())
	r.effect.Init(pixels)
	r.canvas.Set(pixels)
	r.canvas.Draw(0)

	for {
		select {
		case <-r.ctx.Done():
			log.Infof("[%s] Runner stopped", r.name)
			return
		case <-time.After(r.delay):
			log.Debugf("[%s] Rendering: %#v", r.name, r.effect)
			r.effect.Draw(pixels)
			r.canvas.Set(pixels)
			r.canvas.Draw(r.delay)
		}
	}
}

var (
	ErrPresetNotFound   = errors.New("preset not found")
	ErrPresetNotRunning = errors.New("preset not running")
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

// Return the active runner's configuration for this preset if available,
// otherwise return the default configuration for that preset.
func (r *manager) Configuration(preset string) map[string]interface{} {
	buf := &bytes.Buffer{}
	run := r.runners[preset]
	json.NewEncoder(buf).Encode(run.effect)
	values := make(map[string]interface{})
	json.NewDecoder(buf).Decode(&values)
	return values
}

func (r *manager) Configure(preset string, values map[string]interface{}) error {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(values)
	run := r.runners[preset]
	return json.NewDecoder(buf).Decode(run.effect)
}

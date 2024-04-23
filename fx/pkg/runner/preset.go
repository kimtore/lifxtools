package runner

import (
	"context"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/dorkowscy/lifxtool/pkg/config"
	"github.com/dorkowscy/lifxtool/pkg/effects"
)

func initRunners(presets map[string]config.Preset, canvases map[string]canvas.Canvas) (map[string]*runner, error) {
	runners := make(map[string]*runner)
	for preset := range presets {
		run, err := initRunner(presets[preset], canvases[presets[preset].Canvas])
		if err != nil {
			return nil, err
		}
		runners[preset] = run
	}
	return runners, nil
}

func initRunner(preset config.Preset, canvas canvas.Canvas) (*runner, error) {
	eff, err := effects.NewFromConfig(preset.Effect)
	if err != nil {
		return nil, err
	}

	// Prevent effect from running immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return &runner{
		ctx:    ctx,
		cancel: cancel,
		canvas: canvas,
		delay:  preset.Delay,
		effect: eff,
		name:   preset.Name,
	}, nil
}

package effects

import (
	"context"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

type runner struct {
	canvas canvas.Canvas
	cancel context.CancelFunc
	ctx    context.Context
	delay  time.Duration
	effect Effect
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

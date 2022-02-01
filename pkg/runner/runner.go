package runner

import (
	"context"
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

const (
	lifxProcessingTime = time.Millisecond * 10 // approximate latency added to messages
)

func (r *runner) Run() {
	log.Warnf("[%s] RUNNING @ %.3f fps", r.name, 1.0/r.delay.Seconds())

	pixels := make([]colorful.Color, r.canvas.Size())
	r.effect.Init(pixels)
	r.canvas.Set(pixels)
	r.canvas.Draw(0)
	t := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case <-r.ctx.Done():
			log.Warnf("[%s] STOPPED", r.name)
			return
		case <-t.C:
			t.Reset(r.delay)
			log.Debugf("[%s] RENDER %#v", r.name, r.effect)
			r.effect.Draw(pixels)
			r.canvas.Set(pixels)
			r.canvas.Draw(r.delay - lifxProcessingTime)
		}
	}
}

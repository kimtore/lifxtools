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
	initialDelay = 10 * time.Millisecond
)

func (r *runner) Run() {
	log.Warnf("[%s] RUNNING @ %.3f fps", r.name, 1.0/r.delay.Seconds())

	pixels := make([]colorful.Color, r.canvas.Size())
	r.effect.Init(pixels)
	r.canvas.Set(pixels)
	r.canvas.Draw(0)

	t := time.NewTimer(initialDelay)
	deadline := time.Now().Add(r.delay)

	for {
		select {
		case <-r.ctx.Done():
			log.Warnf("[%s] STOPPED", r.name)
			return
		case <-t.C:
			log.Debugf("[%s] RENDER %#v", r.name, r.effect)
			r.effect.Draw(pixels)
			r.canvas.Set(pixels)
			duration := deadline.Sub(time.Now())
			t.Reset(duration)
			r.canvas.Draw(duration)
			deadline = deadline.Add(r.delay)
		}
	}
}

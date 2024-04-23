package canvas

import (
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type aggregate struct {
	canvases []Canvas
	size     int
}

func (a *aggregate) Fill(color colorful.Color) {
	for _, canvas := range a.canvases {
		canvas.Fill(color)
	}
}

func (a *aggregate) Draw(fadeTime time.Duration) {
	for _, canvas := range a.canvases {
		canvas.Draw(fadeTime)
	}
}

func (a *aggregate) Set(pixels []colorful.Color) {
	start := 0
	for _, canvas := range a.canvases {
		ln := canvas.Size()
		canvas.Set(pixels[start : start+ln])
		start += ln
	}
}

func (a *aggregate) Size() int {
	return a.size
}

var _ Canvas = &aggregate{}

func NewAggregate(cvs ...Canvas) Canvas {
	size := 0
	for _, cv := range cvs {
		size += cv.Size()
	}
	return &aggregate{
		canvases: cvs,
		size:     size,
	}
}

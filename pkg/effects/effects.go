package effects

import (
	"github.com/lucasb-eyer/go-colorful"
)

type Effect interface {
	Draw(pixels []colorful.Color)
}

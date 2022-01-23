package canvas

import (
	"time"

	"github.com/dorkowscy/lyslix/lifx"
	"github.com/lucasb-eyer/go-colorful"
)

type Canvas interface {
	Fill(color colorful.Color)
	Draw(fadeTime time.Duration)
	Set(pixels []colorful.Color)
	Pixels() []colorful.Color
}

// Convert a colorful.Color into LIFX's representation of color.
func HBSK(color colorful.Color) lifx.HBSK {
	const full = float64(65535)
	h, s, v := color.Clamped().Hsv()
	return lifx.HBSK{
		Hue:        uint16((h / 360) * full),
		Saturation: uint16(s * full),
		Brightness: uint16(v * full),
		Kelvin:     uint16(6500),
	}
}

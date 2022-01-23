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

func ToHBSK(color colorful.Color) lifx.HBSK {
	const max = float64(65535)
	h, s, v := color.Clamped().Hsv()
	return lifx.HBSK{
		Hue:        uint16((h / 360) * max),
		Saturation: uint16(s * max),
		Brightness: uint16(v * max),
		Kelvin:     uint16(6500),
	}
}

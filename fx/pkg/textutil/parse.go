package textutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

func ParseRGB(input string) (colorful.Color, error) {
	floats, err := stringFloats(input, 3)
	if err != nil {
		return colorful.Color{}, err
	}
	return colorful.LinearRgb(floats[0]/256.0, floats[1]/256.0, floats[2]/256.0), nil
}

func ParseHCL(input string) (colorful.Color, error) {
	if !strings.HasPrefix(input, "hcl(") || !strings.HasSuffix(input, ")") {
		return colorful.Color{}, fmt.Errorf("invalid format: '%s'; expected 'hcl(hue, chroma, luminance)'", input)
	}
	input = strings.TrimLeft(input, "hcl(")
	input = strings.TrimRight(input, ")")

	floats, err := stringFloats(input, 3)
	if err != nil {
		return colorful.Color{}, err
	}

	return colorful.Hcl(floats[0], floats[1], floats[2]), nil
}

func stringFloats(input string, num int) ([]float64, error) {
	floats := make([]float64, 0)
	parts := strings.Split(input, ",")
	if len(parts) != num {
		return nil, fmt.Errorf("string contains %d numbers; expected %d", len(parts), num)
	}
	for _, s := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return nil, err
		}
		floats = append(floats, f)
	}
	return floats, nil
}

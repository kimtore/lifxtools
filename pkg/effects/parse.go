package effects

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type Color struct {
	colorful.Color
}

var _ json.Unmarshaler = &Color{}

func (c *Color) UnmarshalJSON(bytes []byte) error {
	input := ""
	err := json.Unmarshal(bytes, &input)
	if err != nil {
		return err
	}
	color, err := ParseRGB(input)
	if err != nil {
		color, err = ParseHCL(input)
	}
	if err == nil {
		c.Color = color
		return nil
	}
	return err
}

func (c *Color) MarshalJSON() ([]byte, error) {
	r, g, b := c.Clamped().RGB255()
	s := fmt.Sprintf(`"%d,%d,%d"`, r, g, b)
	return []byte(s), nil
}

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

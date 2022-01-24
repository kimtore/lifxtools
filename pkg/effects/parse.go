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
	c.Color, err = ParseHCL(input)
	return err
}

func ParseHCL(input string) (colorful.Color, error) {
	s := strings.TrimLeft(input, "hcl(")
	s = strings.TrimRight(s, ")")
	parts := strings.Split(s, ",")

	if len(parts) != 3 {
		return colorful.Color{}, fmt.Errorf("invalid format: '%s'; expected 'hcl(hue, chroma, luminance)'", input)
	}

	hue, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	chroma, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	luminance, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)

	return colorful.Hcl(hue, chroma, luminance), nil
}

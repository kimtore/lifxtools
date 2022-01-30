package effects

import (
	"encoding/json"
	"fmt"

	"github.com/dorkowscy/lifxtool/pkg/textutil"
	"github.com/lucasb-eyer/go-colorful"
)

type Color struct {
	colorful.Color
}

var _ json.Marshaler = &Color{}
var _ json.Unmarshaler = &Color{}

func (c *Color) UnmarshalJSON(bytes []byte) error {
	input := ""
	err := json.Unmarshal(bytes, &input)
	if err != nil {
		return err
	}
	color, err := textutil.ParseRGB(input)
	if err != nil {
		color, err = textutil.ParseHCL(input)
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

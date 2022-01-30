package textutil

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

func FprintColorHSV(w io.Writer, color colorful.Color) {
	h, s, v := color.Hsv()
	fmt.Fprintf(w, "H=%-9.5f S=%.5f V=%.5f", h, s, v)
	if color.Clamped() != color {
		fmt.Fprintf(w, " [out of gamut]")
	}
	fmt.Fprintf(w, "\n")
}

func FprintColorHCL(w io.Writer, color colorful.Color) {
	h, c, l := color.Hcl()
	fmt.Fprintf(w, "H*=%-9.5f C*=%.5f l*=%.5f", h, c, l)
	if color.Clamped() != color {
		fmt.Fprintf(w, " [out of gamut]")
	}
	fmt.Fprintf(w, "\n")
}

func PrintColorHCL(color colorful.Color) {
	buf := &bytes.Buffer{}
	FprintColorHCL(buf, color)
	log.Debug(buf.String())
}

func PrintColorHSV(color colorful.Color) {
	buf := &bytes.Buffer{}
	FprintColorHSV(buf, color)
	log.Debug(buf.String())
}

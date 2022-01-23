package effects_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/dorkowscy/lifxtool/pkg/effects"
	"github.com/lucasb-eyer/go-colorful"
)

func fprintColorHSV(w io.Writer, color colorful.Color) {
	h, s, v := color.Hsv()
	fmt.Fprintf(w, "H=%-9.5f S=%.5f V=%.5f", h, s, v)
	if color.Clamped() != color {
		fmt.Fprintf(w, " [out of gamut]")
	}
	fmt.Fprintf(w, "\n")
}

func fprintColorHCL(w io.Writer, color colorful.Color) {
	h, c, l := color.Hcl()
	fmt.Fprintf(w, "H*=%-9.5f C*=%.5f l*=%.5f", h, c, l)
	if color.Clamped() != color {
		fmt.Fprintf(w, " [out of gamut]")
	}
	fmt.Fprintf(w, "\n")
}

func printColorHCL(color colorful.Color) {
	fprintColorHCL(os.Stdout, color)
}

func printColorHSV(color colorful.Color) {
	fprintColorHSV(os.Stdout, color)
}

func TestHCLCircle(t *testing.T) {
	colors := effects.HCLCircle(0, 0.1, 0.6, 200)
	for _, color := range colors {
		printColorHSV(color)
	}
}

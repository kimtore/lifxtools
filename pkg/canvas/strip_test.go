package canvas_test

import (
	"testing"
	"time"

	"github.com/dorkowscy/lifxtool/mocks"
	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/lucasb-eyer/go-colorful"
)

// Tests for color zone optimizations.

const size = 10
const min = 0
const max = 9
const fadetime = time.Duration(0)

// Test that pixels don't get sent the same value twice.
func TestStrip_DrawBlankCached(t *testing.T) {
	hbsk := canvas.HBSK(colorful.Color{})
	cli := &mocks.Client{}
	cli.On("SetColorZones", hbsk, uint8(0), uint8(9), fadetime).Return(nil).Once()

	cv := canvas.NewStrip(cli, min, max)
	cv.Draw(fadetime)

	// extra draws should be ignored by cache
	cv.Draw(fadetime)
	cv.Draw(fadetime)

	cli.AssertExpectations(t)
}

// Test that a range of zones in the middle of a black strip
// should be drawn in three separate operations.
func TestStrip_DrawContiguousZone(t *testing.T) {
	black := colorful.Color{}
	color := colorful.Hcl(180, 0.25, 0.25)
	cli := &mocks.Client{}

	cli.On("SetColorZones", canvas.HBSK(black), uint8(0), uint8(3), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.HBSK(color), uint8(4), uint8(8), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.HBSK(black), uint8(9), uint8(9), fadetime).Return(nil).Once()

	cv := canvas.NewStrip(cli, min, max)
	pixels := make([]colorful.Color, cv.Size())

	for i := 4; i <= 8; i++ {
		pixels[i] = color
	}

	cv.Set(pixels)
	cv.Draw(fadetime)

	// Then, test that nulling out those zones result in a draw that affects only those zones.
	cli.On("SetColorZones", canvas.HBSK(black), uint8(4), uint8(8), fadetime).Return(nil).Once()

	cv.Fill(black)
	cv.Draw(fadetime)

	cli.AssertExpectations(t)
}

// Just like DrawContiguousZone, but do the same with two zones, resulting in five operations.
func TestStrip_DrawTwoZones(t *testing.T) {
	black := colorful.Color{}
	color := colorful.Hcl(180, 0.25, 0.25)
	cli := &mocks.Client{}

	cli.On("SetColorZones", canvas.HBSK(black), uint8(0), uint8(3), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.HBSK(color), uint8(4), uint8(4), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.HBSK(black), uint8(5), uint8(6), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.HBSK(color), uint8(7), uint8(7), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.HBSK(black), uint8(8), uint8(9), fadetime).Return(nil).Once()

	cv := canvas.NewStrip(cli, min, max)
	pixels := make([]colorful.Color, cv.Size())

	pixels[4] = color
	pixels[7] = color

	cv.Set(pixels)
	cv.Draw(fadetime)

	cli.AssertExpectations(t)
}

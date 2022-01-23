package canvas_test

import (
	"testing"
	"time"

	"github.com/dorkowscy/lifxtool/mocks"
	"github.com/dorkowscy/lifxtool/pkg/canvas"
	"github.com/lucasb-eyer/go-colorful"
)

const size = 10
const fadetime = time.Duration(0)

func TestStrip_DrawBlankCached(t *testing.T) {
	hbsk := canvas.ToHBSK(colorful.Color{})
	cli := &mocks.Client{}
	cli.On("SetColorZones", hbsk, uint8(0), uint8(9), fadetime).Return(nil).Once()

	cv := canvas.NewStrip(cli, size)
	cv.Draw(fadetime)

	// extra draws should be ignored by cache
	cv.Draw(fadetime)
	cv.Draw(fadetime)
}

func TestStrip_DrawSinglePixel(t *testing.T) {
	black := colorful.Color{}
	color := colorful.Hcl(180, 0.25, 0.25)
	cli := &mocks.Client{}

	cli.On("SetColorZones", canvas.ToHBSK(black), uint8(0), uint8(3), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.ToHBSK(color), uint8(4), uint8(8), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.ToHBSK(black), uint8(9), uint8(9), fadetime).Return(nil).Once()

	cv := canvas.NewStrip(cli, size)
	pixels := cv.Pixels()

	for i := 4; i <= 8; i++ {
		pixels[i] = color
	}

	cv.Set(pixels)
	cv.Draw(fadetime)

	cli.On("SetColorZones", canvas.ToHBSK(black), uint8(4), uint8(8), fadetime).Return(nil).Once()

	cv.Fill(black)
	cv.Draw(fadetime)

	cli.AssertExpectations(t)
}

func TestStrip_DrawDualPixel(t *testing.T) {
	black := colorful.Color{}
	color := colorful.Hcl(180, 0.25, 0.25)
	cli := &mocks.Client{}

	cli.On("SetColorZones", canvas.ToHBSK(black), uint8(0), uint8(3), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.ToHBSK(color), uint8(4), uint8(4), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.ToHBSK(black), uint8(5), uint8(6), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.ToHBSK(color), uint8(7), uint8(7), fadetime).Return(nil).Once()
	cli.On("SetColorZones", canvas.ToHBSK(black), uint8(8), uint8(9), fadetime).Return(nil).Once()

	cv := canvas.NewStrip(cli, size)
	pixels := cv.Pixels()

	pixels[4] = color
	pixels[7] = color

	cv.Set(pixels)
	cv.Draw(fadetime)

	cli.AssertExpectations(t)
}

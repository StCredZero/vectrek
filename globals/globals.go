package globals

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

var (
	WhiteImage = ebiten.NewImage(3, 3)

	// WhiteSubImage is an internal sub image of WhiteImage.
	// Use WhiteSubImage at DrawTriangles instead of WhiteImage in order to avoid bleeding edges.
	WhiteSubImage = WhiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	WhiteImage.Fill(color.White)
}

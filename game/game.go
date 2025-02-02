package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
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

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	Counter int

	AA   bool
	Line bool

	Vertices []ebiten.Vertex
	Indices  []uint16

	Ship *Ship // Player's spaceship
}

// NewShip creates a new ship at the center of the screen
func NewShip() *Ship {
	return &Ship{
		X:     ScreenWidth / 2,
		Y:     ScreenHeight / 2,
		Angle: 0.0,
	}
}

// Ship represents the player's spaceship with position, rotation, and movement
type Ship struct {
	X        float64 // X position on screen
	Y        float64 // Y position on screen
	XV       float64
	YV       float64
	Angle    float64 // rotation Angle in radians
	Velocity float64 // current Velocity
}

type ShipInput struct {
	Left   bool
	Right  bool
	Thrust bool
}

// Ship thrust
const (
	thrustAccel = 0.2
	maxVelocity = 5.0
)

func (s *Ship) Update(input ShipInput) error {
	// Ship rotation (3 degrees per tick)
	if input.Left {
		s.Angle -= 3 * (math.Pi / 180)
	}
	if input.Right {
		s.Angle += 3 * (math.Pi / 180)
	}
	if input.Thrust {
		// Update ship position based on velocity and angle
		s.XV += math.Cos(s.Angle) * thrustAccel
		s.YV += math.Sin(s.Angle) * thrustAccel
	}
	s.X += s.XV
	s.Y += s.YV
	// Wrap around screen edges (toroidal topology)
	if s.X < 0 {
		s.X += ScreenWidth
	} else if s.X >= ScreenWidth {
		s.X -= ScreenWidth
	}
	if s.Y < 0 {
		s.Y += ScreenHeight
	} else if s.Y >= ScreenHeight {
		s.Y -= ScreenHeight
	}
	return nil
}

func (g *Game) drawEbitenText(screen *ebiten.Image, x, y int, aa bool, line bool) {
	var path vector.Path

	// E
	path.MoveTo(20, 20)
	path.LineTo(20, 70)
	path.LineTo(70, 70)
	path.LineTo(70, 60)
	path.LineTo(30, 60)
	path.LineTo(30, 50)
	path.LineTo(70, 50)
	path.LineTo(70, 40)
	path.LineTo(30, 40)
	path.LineTo(30, 30)
	path.LineTo(70, 30)
	path.LineTo(70, 20)
	path.Close()

	// B
	path.MoveTo(80, 20)
	path.LineTo(80, 70)
	path.LineTo(100, 70)
	path.QuadTo(150, 57.5, 100, 45)
	path.QuadTo(150, 32.5, 100, 20)
	path.Close()

	// I
	path.MoveTo(140, 20)
	path.LineTo(140, 70)
	path.LineTo(150, 70)
	path.LineTo(150, 20)
	path.Close()

	// T
	path.MoveTo(160, 20)
	path.LineTo(160, 30)
	path.LineTo(180, 30)
	path.LineTo(180, 70)
	path.LineTo(190, 70)
	path.LineTo(190, 30)
	path.LineTo(210, 30)
	path.LineTo(210, 20)
	path.Close()

	// E
	path.MoveTo(220, 20)
	path.LineTo(220, 70)
	path.LineTo(270, 70)
	path.LineTo(270, 60)
	path.LineTo(230, 60)
	path.LineTo(230, 50)
	path.LineTo(270, 50)
	path.LineTo(270, 40)
	path.LineTo(230, 40)
	path.LineTo(230, 30)
	path.LineTo(270, 30)
	path.LineTo(270, 20)
	path.Close()

	// N
	path.MoveTo(280, 20)
	path.LineTo(280, 70)
	path.LineTo(290, 70)
	path.LineTo(290, 35)
	path.LineTo(320, 70)
	path.LineTo(330, 70)
	path.LineTo(330, 20)
	path.LineTo(320, 20)
	path.LineTo(320, 55)
	path.LineTo(290, 20)
	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 5
		op.LineJoin = vector.LineJoinRound
		g.Vertices, g.Indices = path.AppendVerticesAndIndicesForStroke(g.Vertices[:0], g.Indices[:0], op)
	} else {
		g.Vertices, g.Indices = path.AppendVerticesAndIndicesForFilling(g.Vertices[:0], g.Indices[:0])
	}

	for i := range g.Vertices {
		g.Vertices[i].DstX = (g.Vertices[i].DstX + float32(x))
		g.Vertices[i].DstY = (g.Vertices[i].DstY + float32(y))
		g.Vertices[i].SrcX = 1
		g.Vertices[i].SrcY = 1
		g.Vertices[i].ColorR = 0xdb / float32(0xff)
		g.Vertices[i].ColorG = 0x56 / float32(0xff)
		g.Vertices[i].ColorB = 0x20 / float32(0xff)
		g.Vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa

	// For strokes (AppendVerticesAndIndicesForStroke), FillRuleFillAll and FillRuleNonZero work.
	//
	// For filling (AppendVerticesAndIndicesForFilling), FillRuleNonZero and FillRuleEvenOdd work.
	// FillRuleNonZero and FillRuleEvenOdd differ when rendering a complex polygons with self-intersections and/or holes.
	// See https://en.wikipedia.org/wiki/Nonzero-rule and https://en.wikipedia.org/wiki/Even%E2%80%93odd_rule .
	//
	// For simplicity, this example always uses FillRuleNonZero, whichever strokes or filling is done.
	op.FillRule = ebiten.FillRuleNonZero

	screen.DrawTriangles(g.Vertices, g.Indices, WhiteSubImage, op)
}

func (g *Game) drawArc(screen *ebiten.Image, count int, aa bool, line bool) {
	var path vector.Path

	path.MoveTo(350, 100)
	const cx, cy, r = 450, 100, 70
	theta1 := math.Pi * float64(count) / 180
	x := cx + r*math.Cos(theta1)
	y := cy + r*math.Sin(theta1)
	path.ArcTo(450, 100, float32(x), float32(y), 30)
	path.LineTo(float32(x), float32(y))

	theta2 := math.Pi * float64(count) / 180 / 3
	path.MoveTo(550, 100)
	path.Arc(550, 100, 50, float32(theta1), float32(theta2), vector.Clockwise)
	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 5
		op.LineJoin = vector.LineJoinRound
		g.Vertices, g.Indices = path.AppendVerticesAndIndicesForStroke(g.Vertices[:0], g.Indices[:0], op)
	} else {
		g.Vertices, g.Indices = path.AppendVerticesAndIndicesForFilling(g.Vertices[:0], g.Indices[:0])
	}

	for i := range g.Vertices {
		g.Vertices[i].SrcX = 1
		g.Vertices[i].SrcY = 1
		g.Vertices[i].ColorR = 0x33 / float32(0xff)
		g.Vertices[i].ColorG = 0xcc / float32(0xff)
		g.Vertices[i].ColorB = 0x66 / float32(0xff)
		g.Vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(g.Vertices, g.Indices, WhiteSubImage, op)
}

func (g *Game) Update() error {
	g.Counter++

	// Switch anti-alias.
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.AA = !g.AA
	}

	// Switch lines.
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.Line = !g.Line
	}

	var shipInput ShipInput
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		shipInput.Left = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		shipInput.Right = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		fmt.Println("up")
		shipInput.Thrust = true
	}

	return nil
}

func (g *Game) drawShip(screen *ebiten.Image, aa bool, line bool) {
	var path vector.Path

	// Define ship as a triangle
	length := float32(15.0)
	theta := float32(g.Ship.Angle)

	// Front point
	path.MoveTo(
		float32(g.Ship.X)+length*float32(math.Cos(float64(theta))),
		float32(g.Ship.Y)+length*float32(math.Sin(float64(theta))),
	)

	// Right point (120 degrees from front)
	path.LineTo(
		float32(g.Ship.X)+length*float32(math.Cos(float64(theta)+2.0944)), // 2.0944 rad = 120 deg
		float32(g.Ship.Y)+length*float32(math.Sin(float64(theta)+2.0944)),
	)

	// Left point (-120 degrees from front)
	path.LineTo(
		float32(g.Ship.X)+length*float32(math.Cos(float64(theta)-2.0944)),
		float32(g.Ship.Y)+length*float32(math.Sin(float64(theta)-2.0944)),
	)

	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 2
		op.LineJoin = vector.LineJoinRound
		g.Vertices, g.Indices = path.AppendVerticesAndIndicesForStroke(g.Vertices[:0], g.Indices[:0], op)
	} else {
		g.Vertices, g.Indices = path.AppendVerticesAndIndicesForFilling(g.Vertices[:0], g.Indices[:0])
	}

	for i := range g.Vertices {
		g.Vertices[i].SrcX = 1
		g.Vertices[i].SrcY = 1
		g.Vertices[i].ColorR = 1
		g.Vertices[i].ColorG = 1
		g.Vertices[i].ColorB = 1
		g.Vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(g.Vertices, g.Indices, WhiteSubImage, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	dst := screen

	dst.Fill(color.RGBA{0x00, 0x00, 0x00, 0xff})
	g.drawEbitenText(dst, 0, 50, g.AA, g.Line)
	g.drawArc(dst, g.Counter, g.AA, g.Line)
	g.drawShip(dst, g.AA, g.Line)

	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	msg += "\nPress A to switch anti-alias."
	msg += "\nPress L to switch the fill mode and the line mode."
	msg += "\nUse arrow keys to control the ship."
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

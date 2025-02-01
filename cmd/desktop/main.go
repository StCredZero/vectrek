package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	counter int

	aa   bool
	line bool

	vertices []ebiten.Vertex
	indices  []uint16

	ship *Ship // Player's spaceship
}

// NewShip creates a new ship at the center of the screen
func NewShip() *Ship {
	return &Ship{
		x:     screenWidth / 2,
		y:     screenHeight / 2,
		angle: 0.0,
	}
}

// Ship represents the player's spaceship with position, rotation, and movement
type Ship struct {
	x        float64 // x position on screen
	y        float64 // y position on screen
	angle    float64 // rotation angle in radians
	velocity float64 // current velocity
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
		g.vertices, g.indices = path.AppendVerticesAndIndicesForStroke(g.vertices[:0], g.indices[:0], op)
	} else {
		g.vertices, g.indices = path.AppendVerticesAndIndicesForFilling(g.vertices[:0], g.indices[:0])
	}

	for i := range g.vertices {
		g.vertices[i].DstX = (g.vertices[i].DstX + float32(x))
		g.vertices[i].DstY = (g.vertices[i].DstY + float32(y))
		g.vertices[i].SrcX = 1
		g.vertices[i].SrcY = 1
		g.vertices[i].ColorR = 0xdb / float32(0xff)
		g.vertices[i].ColorG = 0x56 / float32(0xff)
		g.vertices[i].ColorB = 0x20 / float32(0xff)
		g.vertices[i].ColorA = 1
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

	screen.DrawTriangles(g.vertices, g.indices, whiteSubImage, op)
}

func (g *Game) drawEbitenLogo(screen *ebiten.Image, x, y int, aa bool, line bool) {
	const unit = 16

	var path vector.Path

	// TODO: Add curves
	path.MoveTo(0, 4*unit)
	path.LineTo(0, 6*unit)
	path.LineTo(2*unit, 6*unit)
	path.LineTo(2*unit, 5*unit)
	path.LineTo(3*unit, 5*unit)
	path.LineTo(3*unit, 4*unit)
	path.LineTo(4*unit, 4*unit)
	path.LineTo(4*unit, 2*unit)
	path.LineTo(6*unit, 2*unit)
	path.LineTo(6*unit, 1*unit)
	path.LineTo(5*unit, 1*unit)
	path.LineTo(5*unit, 0)
	path.LineTo(4*unit, 0)
	path.LineTo(4*unit, 2*unit)
	path.LineTo(2*unit, 2*unit)
	path.LineTo(2*unit, 3*unit)
	path.LineTo(unit, 3*unit)
	path.LineTo(unit, 4*unit)
	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 5
		op.LineJoin = vector.LineJoinRound
		g.vertices, g.indices = path.AppendVerticesAndIndicesForStroke(g.vertices[:0], g.indices[:0], op)
	} else {
		g.vertices, g.indices = path.AppendVerticesAndIndicesForFilling(g.vertices[:0], g.indices[:0])
	}

	for i := range g.vertices {
		g.vertices[i].DstX = (g.vertices[i].DstX + float32(x))
		g.vertices[i].DstY = (g.vertices[i].DstY + float32(y))
		g.vertices[i].SrcX = 1
		g.vertices[i].SrcY = 1
		g.vertices[i].ColorR = 0xdb / float32(0xff)
		g.vertices[i].ColorG = 0x56 / float32(0xff)
		g.vertices[i].ColorB = 0x20 / float32(0xff)
		g.vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(g.vertices, g.indices, whiteSubImage, op)
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
		g.vertices, g.indices = path.AppendVerticesAndIndicesForStroke(g.vertices[:0], g.indices[:0], op)
	} else {
		g.vertices, g.indices = path.AppendVerticesAndIndicesForFilling(g.vertices[:0], g.indices[:0])
	}

	for i := range g.vertices {
		g.vertices[i].SrcX = 1
		g.vertices[i].SrcY = 1
		g.vertices[i].ColorR = 0x33 / float32(0xff)
		g.vertices[i].ColorG = 0xcc / float32(0xff)
		g.vertices[i].ColorB = 0x66 / float32(0xff)
		g.vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(g.vertices, g.indices, whiteSubImage, op)
}

func maxCounter(index int) int {
	return 128 + (17*index+32)%64
}

func (g *Game) drawWave(screen *ebiten.Image, counter int, aa bool, line bool) {
	var path vector.Path

	const npoints = 8
	indexToPoint := func(i int, counter int) (float32, float32) {
		x, y := float32(i*screenWidth/(npoints-1)), float32(screenHeight/2)
		y += float32(30 * math.Sin(float64(counter)*2*math.Pi/float64(maxCounter(i))))
		return x, y
	}

	for i := 0; i <= npoints; i++ {
		if i == 0 {
			path.MoveTo(indexToPoint(i, counter))
			continue
		}
		cpx0, cpy0 := indexToPoint(i-1, counter)
		x, y := indexToPoint(i, counter)
		cpx1, cpy1 := x, y
		cpx0 += 30
		cpx1 -= 30
		path.CubicTo(cpx0, cpy0, cpx1, cpy1, x, y)
	}
	path.LineTo(screenWidth, screenHeight)
	path.LineTo(0, screenHeight)

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 5
		op.LineJoin = vector.LineJoinRound
		g.vertices, g.indices = path.AppendVerticesAndIndicesForStroke(g.vertices[:0], g.indices[:0], op)
	} else {
		g.vertices, g.indices = path.AppendVerticesAndIndicesForFilling(g.vertices[:0], g.indices[:0])
	}

	for i := range g.vertices {
		g.vertices[i].SrcX = 1
		g.vertices[i].SrcY = 1
		g.vertices[i].ColorR = 0x33 / float32(0xff)
		g.vertices[i].ColorG = 0x66 / float32(0xff)
		g.vertices[i].ColorB = 0xff / float32(0xff)
		g.vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(g.vertices, g.indices, whiteSubImage, op)
}

func (g *Game) Update() error {
	g.counter++

	// Switch anti-alias.
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.aa = !g.aa
	}

	// Switch lines.
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.line = !g.line
	}

	// Ship rotation (3 degrees per frame)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.ship.angle -= 3 * (math.Pi / 180)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.ship.angle += 3 * (math.Pi / 180)
	}

	// Ship thrust
	const (
		thrustAccel = 0.2
		maxVelocity = 5.0
	)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.ship.velocity += thrustAccel
		if g.ship.velocity > maxVelocity {
			g.ship.velocity = maxVelocity
		}
	} else {
		// Apply slight drag when not thrusting
		g.ship.velocity *= 0.99
	}

	// Update ship position based on velocity and angle
	g.ship.x += g.ship.velocity * math.Cos(g.ship.angle)
	g.ship.y += g.ship.velocity * math.Sin(g.ship.angle)

	// Wrap around screen edges (toroidal topology)
	if g.ship.x < 0 {
		g.ship.x += screenWidth
	} else if g.ship.x >= screenWidth {
		g.ship.x -= screenWidth
	}
	if g.ship.y < 0 {
		g.ship.y += screenHeight
	} else if g.ship.y >= screenHeight {
		g.ship.y -= screenHeight
	}

	return nil
}

func (g *Game) drawShip(screen *ebiten.Image, aa bool, line bool) {
	var path vector.Path

	// Define ship as a triangle
	length := float32(15.0)
	theta := float32(g.ship.angle)
	
	// Front point
	path.MoveTo(
		float32(g.ship.x)+length*float32(math.Cos(float64(theta))),
		float32(g.ship.y)+length*float32(math.Sin(float64(theta))),
	)
	
	// Right point (120 degrees from front)
	path.LineTo(
		float32(g.ship.x)+length*float32(math.Cos(float64(theta)+2.0944)), // 2.0944 rad = 120 deg
		float32(g.ship.y)+length*float32(math.Sin(float64(theta)+2.0944)),
	)
	
	// Left point (-120 degrees from front)
	path.LineTo(
		float32(g.ship.x)+length*float32(math.Cos(float64(theta)-2.0944)),
		float32(g.ship.y)+length*float32(math.Sin(float64(theta)-2.0944)),
	)
	
	path.Close()

	if line {
		op := &vector.StrokeOptions{}
		op.Width = 2
		op.LineJoin = vector.LineJoinRound
		g.vertices, g.indices = path.AppendVerticesAndIndicesForStroke(g.vertices[:0], g.indices[:0], op)
	} else {
		g.vertices, g.indices = path.AppendVerticesAndIndicesForFilling(g.vertices[:0], g.indices[:0])
	}

	for i := range g.vertices {
		g.vertices[i].SrcX = 1
		g.vertices[i].SrcY = 1
		g.vertices[i].ColorR = 1
		g.vertices[i].ColorG = 1
		g.vertices[i].ColorB = 1
		g.vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = aa
	op.FillRule = ebiten.FillRuleNonZero
	screen.DrawTriangles(g.vertices, g.indices, whiteSubImage, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	dst := screen

	dst.Fill(color.RGBA{0xe0, 0xe0, 0xe0, 0xff})
	g.drawEbitenText(dst, 0, 50, g.aa, g.line)
	g.drawEbitenLogo(dst, 20, 150, g.aa, g.line)
	g.drawArc(dst, g.counter, g.aa, g.line)
	g.drawWave(dst, g.counter, g.aa, g.line)
	g.drawShip(dst, g.aa, g.line)

	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	msg += "\nPress A to switch anti-alias."
	msg += "\nPress L to switch the fill mode and the line mode."
	msg += "\nUse arrow keys to control the ship."
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		counter: 0,
		ship:    NewShip(),
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Vector (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

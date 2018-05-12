package gobot

import (
	"image"
	"image/color"
	"image/draw"
	"os"
)

// BoardPadding is the padding around the edges of the board
const BoardPadding = 32

// StoneSize is how big each stone cell is in pixels
const StoneSize = 48

// StoneSpacing is how much space to leave around the stones
const StoneSpacing = 2

// BoardPath is the location of the board base image
const BoardPath = "assets/board.png"

var boardImage draw.Image

func init() {
	f, err := os.Open(BoardPath)
	if err != nil {
		panic(err)
	}
	bIm, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	bounds := bIm.Bounds()
	boardImage = image.NewRGBA(bounds)
	draw.Draw(boardImage, bounds, bIm, image.ZP, draw.Src)
}

// Circle is used to draw stones on the Go board
type Circle struct {
	p image.Point
	r int
}

// ColorModel returns an alpha color model
func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

// Bounds gets the bounding rectangle of the circle
func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

// At gets the color of the circle at the given point
func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// Render a board into an image
func Render(board Board) (image.Image, error) {
	im := image.NewRGBA(boardImage.Bounds())
	draw.Draw(im, im.Bounds(), boardImage, image.ZP, draw.Src)
	for i, row := range board {
		for j, stone := range row {
			x := BoardPadding + StoneSize*j + 2*StoneSpacing*j
			y := BoardPadding + StoneSize*i + 2*StoneSpacing*i
			p := image.Point{x, y}
			r := StoneSize / 2
			var src image.Image
			if stone == WhiteStone {
				src = image.White
			} else if stone == BlackStone {
				src = image.Black
			} else {
				src = image.Transparent
			}
			draw.DrawMask(
				im,
				im.Bounds(),
				src,
				image.ZP,
				&Circle{p, r},
				image.ZP,
				draw.Over,
			)
		}
	}
	return im, nil
}

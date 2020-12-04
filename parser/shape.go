package parser

import (
	"github.com/fogleman/gg"
)

type Svg struct {
	Drawables     []drawable
	Width, Height int
}

type drawable interface {
	Draw(dc *gg.Context)
}

type shape struct {
	x, y        float64
	fillColor   string
	strokeColor string
	strokeWidth float64
}

type circle struct {
	shape
	radius float64
}

func (c *circle) Draw(dc *gg.Context) {
	x := c.x
	y := c.y
	r := c.radius

	dc.DrawCircle(x, y, r)
	dc.SetHexColor(c.fillColor)
	dc.Fill()

	dc.DrawCircle(x, y, r)
	dc.SetLineWidth(c.strokeWidth)
	dc.SetHexColor(c.strokeColor)
	dc.Stroke()
}

package rtcore

import (
	"github.com/amartyn1996/go-simple-vecmath"
)

type Camera struct {
	Pos *simple_vecmath.Vector3 
	Dir *simple_vecmath.Vector3
	ViewDistance float64
	ViewWidth float64
	ViewHeight float64
	ImgWidth int
	ImgHeight int
}

package rtcore

import (
	"github.com/amartyn1996/go-simple-vecmath"
)

type Hit struct {
	T float64
	Intersection *simple_vecmath.Vector3
	Normal *simple_vecmath.Vector3
	Surface Surface
}

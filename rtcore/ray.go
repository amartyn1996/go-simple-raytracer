package rtcore

import "github.com/amartyn1996/go-simple-vecmath"

type Ray struct {
	P *simple_vecmath.Vector3
	D *simple_vecmath.Vector3
}

func NewRay() *Ray {
	return &Ray{ simple_vecmath.NewVector3(0,0,0), simple_vecmath.NewVector3(0,0,0) }
}

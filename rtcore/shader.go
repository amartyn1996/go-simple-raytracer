package rtcore

import "github.com/amartyn1996/go-simple-vecmath"

type Shader interface {
	Shade(scene *Scene, hit *Hit) *simple_vecmath.Vector3
}

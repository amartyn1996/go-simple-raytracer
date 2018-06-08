package rtcore

import (
)


type Surface interface {
	Intersect(minT, maxT float64, ray *Ray) *Hit
	GetShader() Shader
}

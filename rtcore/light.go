package rtcore

import (
	"github.com/amartyn1996/go-simple-vecmath"
)

type Light struct {
   	Pos *simple_vecmath.Vector3
	Color *simple_vecmath.Vector3
}

func ComputeLighting(scene *Scene, point *simple_vecmath.Vector3) []*simple_vecmath.Vector3 {

	ret := make([]*simple_vecmath.Vector3, len(scene.Lights) * 2)
	for i , light := range scene.Lights {
		//Create vector L and cast a shadow ray
		PointToLight := light.Pos.Copy().SubtractV(point)
		L := PointToLight.Copy().Normalize()
		AdjustedPoint := point.Copy().ScaleAdd( scene.EPSILON, L )
		shadowCast := scene.CastRay( &Ray{AdjustedPoint,PointToLight} , true)

		//If the shadow ray hit something then don't add this light's color (object is in shadow)
		if (shadowCast != nil && shadowCast.T <= 1) {
			continue
		} else {
			ret[i*2] = L
			ret[i*2+1] = light.Color
		}
	}

	return ret
}

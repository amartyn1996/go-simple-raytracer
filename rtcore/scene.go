package rtcore

import (
    "github.com/amartyn1996/go-simple-vecmath"
    "math"
)

type Scene struct {
	Cam Camera
	Lights []*Light
	Surfaces []Surface
	MinT float64
	MaxT float64
	EPSILON float64
}

func (scene *Scene) CastRay(ray *Ray, shadowRay bool) *Hit {
	minT := math.Inf(1)
	var closest *Hit = nil

	for _,s := range scene.Surfaces {
		hit := s.Intersect(scene.MinT, scene.MaxT, ray)
		if hit != nil && hit.T >= scene.MinT && hit.T < scene.MaxT && hit.T < minT {
			minT = hit.T
			closest = hit
			if shadowRay == true {
    				if (minT <= 1) {
					break;
    				} else {
					closest = nil
    				}
			}
		}
	}

	return closest
}

func (scene *Scene) Trace(ray *Ray) *simple_vecmath.Vector3 {
	returnColor := simple_vecmath.NewVector3(0,0,0)

	closest := scene.CastRay(ray, false)
	
	if closest != nil {
		returnColor.SetV( closest.Surface.GetShader().Shade(scene, closest) )
	}
	return returnColor
}

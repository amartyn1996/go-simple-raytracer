package rtcore

import (
	"github.com/amartyn1996/go-simple-vecmath"
)

type Lambertian struct {
	DiffuseColor *simple_vecmath.Vector3
}

func (s *Lambertian) Shade(scene *Scene, hit *Hit) *simple_vecmath.Vector3 {

    	normal := hit.Normal.Copy().Normalize()
    	lighting := simple_vecmath.NewVector3(0,0,0)

    	lights := ComputeLighting(scene, hit.Intersection) //compute visible lights from light.go
    	for i := 0; i < len(lights); i += 2 {

		L := lights[i]
		lightColor := lights[i+1]

		if L != nil {
			lighting.ScaleAdd( normal.Dot(L), lightColor )
		} 
    	}


    	color := s.DiffuseColor.Copy().ScaleV(lighting)

	return color
}

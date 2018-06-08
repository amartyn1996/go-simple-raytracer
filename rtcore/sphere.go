package rtcore

import (
	"math"
	"github.com/amartyn1996/go-simple-vecmath"
)

type Sphere struct {
	Pos *simple_vecmath.Vector3
	Radius float64
	Shader Shader
}

func (s *Sphere) GetShader() Shader {
	return s.Shader
}

func (s *Sphere) Intersect(minT, maxT float64, ray *Ray) *Hit {

	//line-sphere intersection
	centerToPoint := ray.P.Copy().SubtractV( s.Pos )
	a := ray.D.Length2()
	b := 2 * ray.D.Dot( centerToPoint )
	c := centerToPoint.Length2() - s.Radius

	radical := b*b - 4*a*c

	if radical < 0 { //no intersection
		return nil
	}

	var t float64
	
	if radical == 0 { //one intersection
		t = -b / (2*a)
	} else { //two intersections
    		radical = math.Sqrt(radical)
		t1 := (-b + radical) / (2*a)
		t2 := (-b - radical) / (2*a)

		if !(t2 < minT) && t2 < t1 {
			t = t2
		} else {
			t = t1
		}
	}
	
	if (t < minT && t >= maxT) {
		return nil
	}


	//intersection point and surface normal
	intersection := ray.P.Copy().ScaleAdd(t,ray.D)
	normal := intersection.Copy().SubtractV(s.Pos).Normalize()

	return &Hit{t,intersection,normal,s}
}

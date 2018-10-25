package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
	"strings"
	"github.com/amartyn1996/go-simple-vecmath"
	"./rtcore"
	"./parsers"
)

func main() {
    	startTime := time.Now()
	fmt.Println("Begin.")

	scene, args := initialize()

	img := makeImage(scene)

	writeImage(img, args)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Done. %s\n", elapsedTime)

}

func initialize() (*rtcore.Scene, []string) {
    	startTime := time.Now()
	fmt.Println("Initializing...")

	args := os.Args[1:]

    	s := parsers.Parse(args[0])
	
	elapsedTime := time.Since(startTime)
	fmt.Printf("Initializing Done. %s\n", elapsedTime)

	return s, args
}

func makeImage(scene *rtcore.Scene) image.Image {
    	startTime := time.Now()
	fmt.Println("Creating Image...")

	cam := scene.Cam

	imgWidth := cam.ImgWidth
	imgHeight := cam.ImgHeight

	camPos := cam.Pos.Copy()
	camDir := cam.Dir.Copy()
	dirNorm := camDir.Copy().Normalize()

	//Construct Orthonormal Basis
	up := simple_vecmath.NewVector3(0,1,0)
	w := camDir.Copy().Negate().Normalize()
	u := up.Copy().Cross(w).Normalize()
	v := w.Copy().Cross(u).Normalize()

	//View Plane
	viewDistance := cam.ViewDistance
	viewWidth := cam.ViewWidth
	viewHeight := cam.ViewHeight

	viewPlaneCenter := camPos.Copy().ScaleAdd(viewDistance, dirNorm)
	vpU := u.Copy().Scale(viewWidth)
	vpV := v.Copy().Scale(viewHeight)
	vpEdge := viewPlaneCenter.Copy().ScaleAdd(-0.5,vpU).ScaleAdd(-0.5,vpV)

	fracU := vpU.Copy().Scale(1.0/float64(imgWidth))
	fracV := vpV.Copy().Scale(1.0/float64(imgHeight))
	offsetU := fracU.Copy().Scale(1.0/2.0)
	offsetV := fracV.Copy().Scale(1.0/2.0)

	//Trace Rays
	bounds := image.Rect(0,0,imgWidth, imgHeight)
	img := image.NewNRGBA(bounds)

	tmpU := simple_vecmath.NewVector3(0,0,0)
	tmpV := simple_vecmath.NewVector3(0,0,0)
	vpPoint := simple_vecmath.NewVector3(0,0,0)
	rayDir := simple_vecmath.NewVector3(0,0,0)
	
	ray := rtcore.NewRay()
	for i := 0; i < imgHeight; i++ {
		for j := 0; j < imgWidth; j++ {

			tmpU.SetV(fracU).Scale(float64(j)).AddV(offsetU)
			tmpV.SetV(fracV).Scale(float64(i)).AddV(offsetV)
			vpPoint.SetV(vpEdge).AddV(tmpU).AddV(tmpV)
			rayDir.SetV(vpPoint).SubtractV(camPos).Normalize()

			ray.P.SetV(camPos)
			ray.D.SetV(rayDir)

			colorVec := scene.Trace(ray)

			//gamma correction
			inverseGamma := 1 / 2.2
			colorVec.X = math.Pow(colorVec.X, inverseGamma)
			colorVec.Y = math.Pow(colorVec.Y, inverseGamma)
			colorVec.Z = math.Pow(colorVec.Z, inverseGamma)
			
			colorVec.Clamp(0.0,1.0)
			colorVec.Scale(255)

			var red uint8 = uint8(round(colorVec.X))
			var green uint8 = uint8(round(colorVec.Y))
			var blue uint8 = uint8(round(colorVec.Z))
			color := color.NRGBA{ red, green, blue, uint8(255) }
			img.Set(j,imgHeight-1 - i,color)
		}
	}
	
	elapsedTime := time.Since(startTime)
	fmt.Printf("Image Creation Done. %s\n", elapsedTime)
	return img
}

func writeImage(img image.Image, args []string) {
    	startTime := time.Now()
	fmt.Println("Writing Image to Disk...")

	split := strings.Split(args[0],"/")
	imageName := split[ len(split) - 1 ]
	f, err := os.Create("./scenes/" + imageName + ".png")
	if (err != nil) {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, img)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Done Writing Image. %s\n", elapsedTime)
}


func round(x float64) float64 {
	return math.Floor( x + 0.5 )
}

# Go Simple Raytracer

A simple golang raytracer. It is very much a work in progress and there is a lot that needs to be done before I am happy with it.

Here are a few things I want this project to have when I am done with it:

Shading:
- Phong Shading
- Reflections
- Glass (Refraction)
- Physically Based Shading?

Surfaces:
- Axis-Aligned Box
- Triangle
- Triangle Mesh

Misc:
- Acceleration Structure (Bounding Volume Hierarchy?)
- Transformations
- Concurrency
- Scene XML Parser Generator (Either from a grammar or from the scene file itself)

## Installation and Running the Program

This should do it:

```
go get github.com/amartyn1996/go-simple-raytracer

cd $GOPATH/src/github.com/amartyn1996/go-simple-raytracer/

make

./raytracer scenes/*SCENE_NAME_HERE*
```



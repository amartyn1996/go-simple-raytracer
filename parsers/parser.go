package parsers

import (
	"os"
	//"fmt"
	"bufio"
	"strconv"
	"github.com/amartyn1996/go-simple-vecmath"
	"../rtcore"
)



/*
 * The grammar this parser implements:
 *
 * RaytracerFile -> Shaders Scene
 *
 * Shaders -> SHADERS_TOK MoreShaders SHADERS_END_TOK
 *
 * MoreShaders -> Shader MoreShaders | LAMBDA 
 *
 * Shader -> SHADER_TOK ShaderName ShaderType SHADER_END_TOK
 *
 * ShaderName -> NAME_TOK STRING NAME_END_TOK
 *
 * ShaderType -> Lambertian
 *
 * Lambertian -> LAMBERTIAN_TOK DiffuseColor LAMBERTIAN_END_TOK
 *
 * DiffuseColor -> DIFFUSECOLOR_TOK FLOAT FLOAT FLOAT DIFFUSECOLOR_END_TOK
 *
 * Scene -> SCENE_TOK Camera Lights Surfaces MinT MaxT Epsilon SCENE_END_TOK
 *
 * Camera -> CAMERA_TOK Pos Dir ViewDistance ViewWidth ViewHeight ImgWidth ImgHeight CAMERA_END_TOK
 *
 * Pos -> POS_TOK FLOAT FLOAT FLOAT POS_END_TOK
 *
 * Dir -> DIR_TOK FLOAT FLOAT FLOAT DIR_END_TOK
 *
 * ViewDistance -> VIEWDISTANCE_TOK FLOAT VIEWDISTANCE_END_TOK
 *
 * ViewWidth -> VIEWWIDTH_TOK FLOAT VIEWWIDTH_END_TOK
 *
 * ViewHeight -> VIEWHEIGHT_TOK FLOAT VIEWHEIGHT_END_TOK
 *
 * ImgWidth -> IMGWIDTH_TOK INT IMGWIDTH_END_TOK
 *
 * ImgHeight -> IMGHEIGHT_TOK INT IMGHEIGHT_END_TOK
 *
 * Lights -> LIGHTS_TOK MoreLights LIGHTS_END_TOK
 *
 * MoreLights -> Light MoreLights | LAMBDA
 *
 * Light -> LIGHT_TOK Pos Color LIGHT_END_TOK
 *
 * Color -> COLOR_TOK FLOAT FLOAT FLOAT COLOR_END_TOK
 *
 * Surfaces -> SURFACES_TOK MoreSurfaces SURFACES_END_TOK
 *
 * MoreSurfaces -> Surface MoreSurfaces | LAMBDA
 *
 * Surface -> Sphere
 *
 * Sphere -> SPHERE_TOK Pos Radius ShaderName SPHERE_END_TOK
 *
 * Radius -> RADIUS_TOK FLOAT RADIUS_END_TOK
 *
 * ShaderName -> SHADERNAME_TOK STRING SHADERNAME_END_TOK
 *
 * MinT -> MINT_TOK FLOAT MINT_END_TOK
 *
 * MaxT -> MAXT_TOK FLOAT MAXT_END_TOK
 *
 * Epsilon -> EPSILON_TOK FLOAT EPSILON_END_TOK
*/

const SHADERS_TOK string = "<Shaders>"
const SHADERS_END_TOK string = "</Shaders>"
const SHADER_TOK string = "<Shader>"
const SHADER_END_TOK string = "</Shader>"
const LAMBERTIAN_TOK string = "<Lambertian>"
const LAMBERTIAN_END_TOK string = "</Lambertian>"
const DIFFUSECOLOR_TOK string = "<DiffuseColor>"
const DIFFUSECOLOR_END_TOK string = "</DiffuseColor>"
const SCENE_TOK string = "<Scene>"
const SCENE_END_TOK string = "</Scene>"
const CAMERA_TOK string = "<Camera>"
const CAMERA_END_TOK string = "</Camera>"
const POS_TOK string = "<Pos>"
const POS_END_TOK string = "</Pos>"
const DIR_TOK string = "<Dir>"
const DIR_END_TOK string = "</Dir>"
const VIEWDISTANCE_TOK string = "<ViewDistance>"
const VIEWDISTANCE_END_TOK string = "</ViewDistance>"
const VIEWWIDTH_TOK string = "<ViewWidth>"
const VIEWWIDTH_END_TOK string = "</ViewWidth>"
const VIEWHEIGHT_TOK string = "<ViewHeight>"
const VIEWHEIGHT_END_TOK string = "</ViewHeight>"
const IMGWIDTH_TOK string = "<ImgWidth>"
const IMGWIDTH_END_TOK string = "</ImgWidth>"
const IMGHEIGHT_TOK string = "<ImgHeight>"
const IMGHEIGHT_END_TOK string = "</ImgHeight>"
const LIGHTS_TOK string = "<Lights>"
const LIGHTS_END_TOK string = "</Lights>"
const LIGHT_TOK string = "<Light>"
const LIGHT_END_TOK string = "</Light>"
const COLOR_TOK string = "<Color>"
const COLOR_END_TOK string = "</Color>"
const SURFACES_TOK string = "<Surfaces>"
const SURFACES_END_TOK string = "</Surfaces>"
const SPHERE_TOK string = "<Sphere>"
const SPHERE_END_TOK string = "</Sphere>"
const RADIUS_TOK string = "<Radius>"
const RADIUS_END_TOK string = "</Radius>"
const SHADERNAME_TOK string = "<ShaderName>"
const SHADERNAME_END_TOK string = "</ShaderName>"
const MINT_TOK string = "<MinT>"
const MINT_END_TOK string = "</MinT>"
const MAXT_TOK string = "<MaxT>"
const MAXT_END_TOK string = "</MaxT>"
const EPSILON_TOK string = "<Epsilon>"
const EPSILON_END_TOK string = "</Epsilon>"

type parser struct {
	shaders map[string]rtcore.Shader
	scene *rtcore.Scene
}

type lexer struct {
    	TOK string
	scanner *bufio.Scanner
}

var parse parser = parser{}

var lex lexer = lexer{}

func strToFloat(num string) float64 {
	val, err := strconv.ParseFloat(num, 64)
	if err != nil {
		panic(err)
	}
	return val
}

func (lexer *lexer) nextToken() {
	lexer.scanner.Scan()
	lexer.TOK = lexer.scanner.Text()
}

func Parse(path string) *rtcore.Scene {
	xmlFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()


	//Initialize lexer
	lex.scanner = bufio.NewScanner(xmlFile)
	lex.scanner.Split(bufio.ScanWords)
	lex.nextToken()
	
	//Initialize parser
	parse.shaders = make( map[string]rtcore.Shader )
	parse.scene = &rtcore.Scene{ rtcore.Camera{}, make([]*rtcore.Light,0), make([]rtcore.Surface,0), 0, 0, 0 }


	//Begin the parsing
	raytracerFile()

	return parse.scene
}

func match(s1 string) {
	if s1 == lex.TOK && lex.TOK != "" {
		lex.nextToken()
	} else {
		error()
	}
}

func error() {
	panic("Parse Failed. Current Token: \"" + lex.TOK + "\"")
}

func raytracerFile() {
	shaders()
	scene()
}

func shaders() {
	if lex.TOK == SHADERS_TOK {
		match(SHADERS_TOK)
		moreShaders()
		match(SHADERS_END_TOK)
	} else {
		error()
	}
}

func moreShaders() {
	if lex.TOK == SHADER_TOK {
		shader()
		moreShaders()
	}
}

func shader() {
	if lex.TOK == SHADER_TOK {
		match(SHADER_TOK)
		name := shaderName()
		stype := shaderType()
		match(SHADER_END_TOK)
		parse.shaders[name] = stype
	} else {
		error()
	}
}

func shaderName() string {
    	match(SHADERNAME_TOK)
	name := lex.TOK
	match(name)
	match(SHADERNAME_END_TOK)
	return name
}

func shaderType() rtcore.Shader {
 	var s rtcore.Shader = nil
	if lex.TOK == LAMBERTIAN_TOK {
		s = lambertian()
	} else {
		error()
	}
	return s
}

func lambertian() rtcore.Shader {
    	var color *simple_vecmath.Vector3 = nil
	if lex.TOK == LAMBERTIAN_TOK {
		match(LAMBERTIAN_TOK)
		color = diffuseColor()
		match(LAMBERTIAN_END_TOK)
	} else {
		error()
	}
	return &rtcore.Lambertian{color}
}

func diffuseColor() *simple_vecmath.Vector3 {
    	match(DIFFUSECOLOR_TOK)
	color := simple_vecmath.NewVector3(0,0,0)
	color.X = strToFloat( lex.TOK )
	match(lex.TOK)
	color.Y = strToFloat( lex.TOK )
	match(lex.TOK)
	color.Z = strToFloat( lex.TOK )
	match(lex.TOK)
	match(DIFFUSECOLOR_END_TOK)
	return color
}

func scene() {
	if lex.TOK == SCENE_TOK {
		match(SCENE_TOK)
		camera()
		lights()
		surfaces()
		minT()
		maxT()
		epsilon()
		match(SCENE_END_TOK)
	} else {
		error()
	}
}

func camera() {
	if lex.TOK == CAMERA_TOK {
    		match(CAMERA_TOK)
		parse.scene.Cam.Pos = pos()
		parse.scene.Cam.Dir = dir()
		viewDistance()
		viewWidth()
		viewHeight()
		imgWidth()
		imgHeight()
		match(CAMERA_END_TOK)
	} else {
		error()
	}
}

func pos() *simple_vecmath.Vector3 {
	position := simple_vecmath.NewVector3(0,0,0)
	if lex.TOK == POS_TOK {
		match(POS_TOK)
		position.X = strToFloat( lex.TOK )
		match( lex.TOK )
		position.Y = strToFloat( lex.TOK )
		match( lex.TOK )
		position.Z = strToFloat( lex.TOK )
		match( lex.TOK )
		match(POS_END_TOK)
	} else {
		error()
	}
	return position
}

func dir() *simple_vecmath.Vector3 {
	direction := simple_vecmath.NewVector3(0,0,0)
	if lex.TOK == DIR_TOK {
		match(DIR_TOK)
		direction.X = strToFloat( lex.TOK )
		match( lex.TOK )
		direction.Y = strToFloat( lex.TOK )
		match( lex.TOK )
		direction.Z = strToFloat( lex.TOK )
		match( lex.TOK )
		match(DIR_END_TOK)
	} else {
		error()
	}
	return direction
}

func viewDistance() {
	if lex.TOK == VIEWDISTANCE_TOK {
		match(VIEWDISTANCE_TOK)
		parse.scene.Cam.ViewDistance = strToFloat( lex.TOK )
		match( lex.TOK )
		match(VIEWDISTANCE_END_TOK)
	} else {
		error()
	}
}

func viewWidth() {
	if lex.TOK == VIEWWIDTH_TOK {
		match(VIEWWIDTH_TOK)
		parse.scene.Cam.ViewWidth = strToFloat( lex.TOK )
		match( lex.TOK )
		match(VIEWWIDTH_END_TOK)
	} else {
		error()
	}
}

func viewHeight() {
	if lex.TOK == VIEWHEIGHT_TOK {
		match(VIEWHEIGHT_TOK)
		parse.scene.Cam.ViewHeight = strToFloat( lex.TOK )
		match( lex.TOK )
		match(VIEWHEIGHT_END_TOK)
	} else {
		error()
	}
}

func imgWidth() {
	if lex.TOK == IMGWIDTH_TOK {
		match(IMGWIDTH_TOK)
		parse.scene.Cam.ImgWidth = int( strToFloat( lex.TOK ) )
		match( lex.TOK )
		match(IMGWIDTH_END_TOK)
	} else {
		error()
	}
}

func imgHeight() {
	if lex.TOK == IMGHEIGHT_TOK {
		match(IMGHEIGHT_TOK)
		parse.scene.Cam.ImgHeight = int( strToFloat( lex.TOK ) )
		match( lex.TOK )
		match(IMGHEIGHT_END_TOK)
	} else {
		error()
	}
}

func lights() {
	if lex.TOK == LIGHTS_TOK {
		match(LIGHTS_TOK)
		moreLights()
		match(LIGHTS_END_TOK)
	} else {
		error()
	}
}

func moreLights() {
	if lex.TOK == LIGHT_TOK {
		light()
		moreLights()
	}
}

func light() {
	if lex.TOK == LIGHT_TOK {
		match(LIGHT_TOK)
		light := &rtcore.Light{}
		light.Pos = pos()
		light.Color = colour()
		match(LIGHT_END_TOK)
		parse.scene.Lights = append(parse.scene.Lights, light)
	} else {
		error()
	}
}

func colour() *simple_vecmath.Vector3 { //Changed name from color to colour because of a name conflict
	color := simple_vecmath.NewVector3(0,0,0)
	if lex.TOK == COLOR_TOK {
		match(COLOR_TOK)
		color.X = strToFloat( lex.TOK )
		match( lex.TOK )
		color.Y = strToFloat( lex.TOK )
		match( lex.TOK )
		color.Z = strToFloat( lex.TOK )
		match( lex.TOK )
		match(COLOR_END_TOK)
	} else {
		error()
	}
	return color
}

func surfaces() {
    	if lex.TOK == SURFACES_TOK {
		match(SURFACES_TOK)
		moreSurfaces()
		match(SURFACES_END_TOK)
    	} else {
		error()
    	}
}

func moreSurfaces() {
	if lex.TOK == SPHERE_TOK {
		parse.scene.Surfaces = append(parse.scene.Surfaces, sphere() )
		moreSurfaces()
	}
}

func sphere() rtcore.Surface {
    	surf := &rtcore.Sphere{}
	if lex.TOK == SPHERE_TOK {
		match(SPHERE_TOK)
		surf.Pos = pos()
		surf.Radius = radius()
		surf.Shader = parse.shaders[ shaderName() ]
		match(SPHERE_END_TOK)
	} else {
		error()
	}
	return surf
}

func radius() float64 {
	rad := 0.0
	if lex.TOK == RADIUS_TOK {
		match(RADIUS_TOK)
		rad = strToFloat( lex.TOK )
		match( lex.TOK )
		match(RADIUS_END_TOK)
	} else {
		error()
	}
	return rad
}

func minT() {
	if lex.TOK == MINT_TOK {
		match(MINT_TOK)
		parse.scene.MinT = strToFloat( lex.TOK )
		match( lex.TOK )
		match(MINT_END_TOK)
	} else {
		error()
	}
}

func maxT() {
	if lex.TOK == MAXT_TOK {
		match(MAXT_TOK)
		parse.scene.MaxT = strToFloat( lex.TOK )
		match( lex.TOK )
		match(MAXT_END_TOK)
	} else {
		error()
	}
}

func epsilon() {
	if lex.TOK == EPSILON_TOK {
		match(EPSILON_TOK)
		parse.scene.EPSILON = strToFloat( lex.TOK )
		match( lex.TOK )
		match(EPSILON_END_TOK)
	} else {
		error()
	}
}

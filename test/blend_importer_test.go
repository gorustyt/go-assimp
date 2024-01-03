package test

import (
	"assimp"
	"assimp/common"
	"math"
	"testing"
)

func Test0(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/box.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/box_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test1(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/4Cubes4Mats_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/4Cubes4Mats_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}
func Test2(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/blender_269_regress1.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/blender_269_regress1_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test3(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test4(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_250.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_250_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test5(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_250_Compressed.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_250_Compressed_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test6(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_262.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_262_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test7(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_269.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_269_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test8(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_271.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_271_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test9(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_276.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/BlenderDefault_276_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test10(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/CubeHierarchy_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/CubeHierarchy_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test11(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/HUMAN.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/HUMAN_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test12(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/MirroredCube_252.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/MirroredCube_252_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test13(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/NoisyTexturedCube_VoronoiGlob_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/NoisyTexturedCube_VoronoiGlob_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test14(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/SmoothVsSolidCube_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/SmoothVsSolidCube_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test15(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/Suzanne_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/Suzanne_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test16(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/SuzanneSubdiv_252.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/SuzanneSubdiv_252_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
	// check approximate shape by averaging together all vertices
	Assert(t, len(p.Meshes) == 1)
	var vertexAvg = common.NewAiVector3D3(0.0, 0.0, 0.0)
	for i := 0; i < len(p.Meshes); i++ {
		mesh := p.Meshes[i]
		Assert(t, mesh != nil)

		invVertexCount := 1.0 / float32(len(mesh.Vertices))
		for j := 0; j < len(mesh.Vertices); j++ {
			vertexAvg = vertexAvg.Add(mesh.Vertices[j].Mul(invVertexCount))
		}
	}

	// must not be inf or nan
	Assert(t, math.IsInf(float64(vertexAvg.X), 1) ||
		math.IsInf(float64(vertexAvg.X), -1) || math.IsNaN(float64(vertexAvg.X)))
	Assert(t, math.IsInf(float64(vertexAvg.Y), 1) ||
		math.IsInf(float64(vertexAvg.Y), -1) || math.IsNaN(float64(vertexAvg.Y)))
	Assert(t, math.IsInf(float64(vertexAvg.Z), 1) ||
		math.IsInf(float64(vertexAvg.Z), -1) || math.IsNaN(float64(vertexAvg.Z)))
	AssertFloatEqual(t, float64(vertexAvg.X), 6.4022515289252624e-08, 0.0001)
	AssertFloatEqual(t, float64(vertexAvg.Y), 0.060569953173398972, 0.0001)
	AssertFloatEqual(t, float64(vertexAvg.Z), 0.31429031491279602, 0.0001)
}

func Test17(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/TexturedCube_ImageGlob_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/TexturedCube_ImageGlob_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test18(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/TexturedPlane_ImageUv_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/TexturedPlane_ImageUv_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test19(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/TexturedPlane_ImageUvPacked_248.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/TexturedPlane_ImageUvPacked_248_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test20(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/TorusLightsCams_250_compressed.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/TorusLightsCams_250_compressed_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test21(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/BLEND/yxa_1.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/BLEND/yxa_1_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test22(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_nonbsd_data/BLEND/Bob.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_nonbsd_data/BLEND/Bob_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func Test23(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_nonbsd_data/BLEND/fleurOptonl.blend")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_nonbsd_data/BLEND/fleurOptonl_blend.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

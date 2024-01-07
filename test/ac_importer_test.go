package test

import (
	"github.com/gorustyt/go-assimp"
	"github.com/gorustyt/go-assimp/common"
	"math"
	"testing"
)

func TestAcImporter(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/closedLine.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/closedLine_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter1(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/nosurfaces.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/nosurfaces_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter2(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/openLine.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/openLine_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter3(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/sample_subdiv.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/sample_subdiv_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))

	// check approximate shape by averaging together all vertices
	Assert(t, len(p.Meshes) == 1)
	vertexAvg := common.NewAiVector3D3(0.0, 0.0, 0.0)
	for i := 0; i < len(p.Meshes); i++ {
		mesh := p.Meshes[i]
		Assert(t, mesh != nil)

		invVertexCount := 1.0 / len(mesh.Vertices)
		for j := 0; j < len(mesh.Vertices); j++ {
			vertexAvg = vertexAvg.Add(mesh.Vertices[j].Mul(float32(invVertexCount)))
		}
	}

	// must not be inf or nan
	common.AiAssert(math.IsInf(float64(vertexAvg.X), 1) || math.IsInf(float64(vertexAvg.X), -1))
	common.AiAssert(math.IsInf(float64(vertexAvg.Y), 1) || math.IsInf(float64(vertexAvg.Y), -1))
	common.AiAssert(math.IsInf(float64(vertexAvg.Z), 1) || math.IsInf(float64(vertexAvg.Z), -1))
	common.AiAssert(math.Abs(float64(vertexAvg.X)-0.079997420310974121) < 0.0001)
	common.AiAssert(math.Abs(float64(vertexAvg.Y-0.099498569965362549)) < 0.0001)
	common.AiAssert(math.Abs(float64(vertexAvg.Z-(-0.10344827175140381))) < 0.0001)
}

func TestAcImporter4(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/SphereWithLight.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/SphereWithLight_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter5(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/SphereWithLight_UTF8BOM.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/SphereWithLight_UTF8BOM_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter6(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/SphereWithLight_UTF16LE.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/SphereWithLight_UTF16LE_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))

}

func TestAcImporter7(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/SphereWithLightUvScaling4X.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/SphereWithLightUvScaling4X_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter8(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/TestFormatDetection")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/TestFormatDetection.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter9(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/Wuson.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/Wuson_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter10(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/Wuson.acc")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/Wuson_acc.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter11(t *testing.T) {
	p, err := assimp.ParseFile("../example/example_data/AC/TestFormatDetection")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example/example_data/AC/TestFormatDetection.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

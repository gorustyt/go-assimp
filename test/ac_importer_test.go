package test

import (
	"assimp"
	"testing"
)

func TestAcImporter(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/closedLine.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/closedLine_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter1(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/nosurfaces.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/nosurfaces_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter2(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/openLine.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/openLine_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter3(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/sample_subdiv.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/sample_subdiv_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter4(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/SphereWithLight.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/SphereWithLight_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter5(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/SphereWithLight_UTF8BOM.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/SphereWithLight_UTF8BOM_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter6(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/SphereWithLight_UTF16LE.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/SphereWithLight_UTF16LE_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter7(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/SphereWithLightUvScaling4X.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/SphereWithLightUvScaling4X_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter8(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/TestFormatDetection")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/TestFormatDetection.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter9(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/Wuson.ac")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/Wuson_ac.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

func TestAcImporter10(t *testing.T) {
	p, err := assimp.ParseFile("../example_data/example_data/AC/Wuson.acc")
	AssertError(t, err)
	p1, err := assimp.ParseFile("../example_data/BLEND/Wuson_acc.assbin")
	AssertError(t, err)
	Assert(t, DeepEqual(p, p1))
}

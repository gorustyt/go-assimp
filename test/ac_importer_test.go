package test

import (
	"assimp/common/logger"
	"assimp/driver"
	"go.uber.org/zap"
	"testing"
)

func TestAcImporter(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/closedLine.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter1(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/nosurfaces.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter2(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/openLine.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter3(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/sample_subdiv.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter4(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/SphereWithLight.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter5(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/SphereWithLight_UTF8BOM.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter6(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/SphereWithLight_UTF16LE.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter7(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/SphereWithLightUvScaling4X.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter8(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/TestFormatDetection", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter9(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/Wuson.ac", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

func TestAcImporter10(t *testing.T) {
	im := driver.NewImporter()
	p, err := im.ReadFile("../example_data/exmaple_data/AC/Wuson.acc", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	Assert(t, p != nil)
}

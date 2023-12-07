package AC

import (
	"assimp/common"
	"assimp/core"
	"log"
)

var (
	desc = core.AiImporterDesc{
		"AC3D Importer",
		"",
		"",
		"",
		core.AiImporterFlags_SupportTextFlavour,
		0,
		0,
		0,
		0,
		"ac acc ac3d",
	}
)

func NewAC3DImporter(reader *common.AiReader) *AC3DImporter {
	im := &AC3DImporter{}
	im.BaseImporter.Reader = reader
	return im
}

func (ac *AC3DImporter) CanRead(pFile string, checkSig bool) {

}

func (ac *AC3DImporter) LoadObjectSection(objects []*Object) {

}
func (ac *AC3DImporter) ConvertMaterial(object *Object,
	matSrc *Material,
	matDest *core.AiMaterial) {
}
func (ac *AC3DImporter) ConvertObjectSection(object *Object,
	meshes []core.AiMesh,
	outMaterials []*core.AiMaterial,
	materials []*Material,
	parent *core.AiNode) core.AiNode {

}

func (ac *AC3DImporter) Read() error {
	res, err := ac.Reader.NextKeyString("MATERIAL", 1)
	var mat Material
	if err == nil {
		mat.rgb, err = ac.Reader.NextKeyFloat64("rgb", 3)
		if err != nil {
			return err
		}
		mat.amb, err = ac.Reader.NextKeyFloat64("amb", 3)
		if err != nil {
			return err
		}
		mat.emis, err = ac.Reader.NextKeyFloat64("emis", 3)
		if err != nil {
			return err
		}
		mat.spec, err = ac.Reader.NextKeyFloat64("spec", 3)
		if err != nil {
			return err
		}
		mat.shin, err = ac.Reader.NextKeyFloat64("shi", 1)
		if err != nil {
			return err
		}
		mat.trans, err = ac.Reader.NextKeyFloat64("trans", 1)
		if err != nil {
			return err
		}
		ac.Reader.NextLine()
	} else {
		log.Println("AC3DImporter", err)
	}

}

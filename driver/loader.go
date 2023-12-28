package driver

import (
	"assimp/common/config"
	"assimp/common/logger"
	"assimp/core"
	"assimp/driver/AC"
	"assimp/driver/BLEND"
	"assimp/driver/assetBin"
	"assimp/driver/base/iassimp"
	"assimp/driver/post_processing"
	"assimp/driver/pre_processing"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var (
	loader sync.Map
)

type LoaderCons func(data []byte) (iassimp.Loader, error)

func init() {
	RegisterLoader(ASSETBIN.NewAssBinImporter, ASSETBIN.Desc)
	RegisterLoader(AC.NewAC3DImporter, AC.Desc)
	RegisterLoader(BLEND.NewBlenderImporter, BLEND.Desc)
}

func RegisterLoader(l LoaderCons, desc core.AiImporterDesc) {
	for _, v := range desc.FileExtensions {
		vsi, ok := loader.Load(v)
		var vs []LoaderCons
		if ok {
			vs = vsi.([]LoaderCons)
		}
		vs = append(vs, l)
		loader.Store(v, vs)
	}
}

type importer struct {
	PostProcessingSteps []iassimp.PostProcessing
	bExtraVerbose       bool
}

func NewImporter() iassimp.Importer {
	return &importer{}
}
func (im *importer) ApplyPostProcessing(pScene *core.AiScene, pFlags int) {
	if pFlags == 0 || pScene == nil {
		return
	}
	logger.Info("Entering post processing pipeline")
	// The ValidateDS process plays an exceptional role. It isn't contained in the global
	// list of post-processing steps, so we need to call it manually.
	var ds post_processing.ValidateDSProcess
	ds.Execute(pScene)
	for _, v := range im.PostProcessingSteps {
		if v.IsActive(pFlags) {
			v.Execute(pScene)
		}
	}
	ds = post_processing.ValidateDSProcess{}
	ds.Execute(pScene)
	logger.Info("Leaving post processing pipeline")
}

func (im *importer) ReadFile(path string, pFlags int) (s *core.AiScene, err error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(path)
	vsi, ok := loader.Load(ext)
	var l iassimp.Loader
	try := func(vsi interface{}) bool {
		for _, v := range vsi.([]LoaderCons) {
			if err != nil {
				return false
			}
			ls, ierr := v(data)
			if ierr != nil {
				err = ierr
				return false
			}
			if ls.CanRead(true) {
				l = ls
				return true
			}
		}
		return false
	}
	if ok {
		try(vsi)
	} else {
		loader.Range(func(key, value any) bool {
			return !try(value)
		})
	}
	if l == nil {
		return nil, errors.New("invalid format")
	}
	res := &core.AiScene{}
	cfg := config.NewConfig()
	l.InitConfig(cfg)
	err = l.Read(res)
	if err != nil {
		return res, err
	}
	pre := pre_processing.NewScenePreprocessor(res)
	pre.ProcessScene()
	im.ApplyPostProcessing(res, pFlags&int(^iassimp.AiProcess_ValidateDataStructure))
	return res, nil
}

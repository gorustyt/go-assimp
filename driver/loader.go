package driver

import (
	"assimp/common/reader"
	"assimp/core"
	"assimp/driver/AC"
	"assimp/driver/base/iassimp"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var (
	loader sync.Map
)

type LoaderCons func(aiReader *reader.AiReader) iassimp.Loader

func init() {
	RegisterLoader(AC.NewAC3DImporter, AC.Desc)
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
}

func NewImporter() iassimp.Importer {
	return &importer{}
}

func (im *importer) ReadFile(path string) (*core.AiScene, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	r, err := reader.NewReader(data)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(path)
	vsi, ok := loader.Load(ext)
	var l iassimp.Loader
	try := func(vsi interface{}) bool {
		for _, v := range vsi.([]LoaderCons) {
			ls := v(r)
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
	return res, l.Read(res)
}

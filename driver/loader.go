package driver

import (
	"assimp/common/reader"
	"assimp/core"
	"assimp/driver/AC"
	"assimp/driver/base/iassimp"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	loader sync.Map
)

func RegisterLoader(name string, l iassimp.Loader) {
	_, ok := loader.Load(name)
	if ok {
		panic(fmt.Sprintf("name:%v loader has exist!", name))
	}
	loader.Store(name, l)
}

func GetLoader() {

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

	l := AC.NewAC3DImporter(r)
	if !l.CanRead(true) {
		return nil, errors.New("invalid format")
	}
	res := &core.AiScene{}
	return res, l.Read(res)
}

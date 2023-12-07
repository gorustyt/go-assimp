package driver

import (
	"assimp/common"
	"assimp/core"
	"fmt"
	"os"
	"sync"
)

var (
	loader sync.Map
)

func RegisterLoader(name string, l Loader) {
	_, ok := loader.Load(name)
	if ok {
		panic(fmt.Sprintf("name:%v loader has exist!", name))
	}
	loader.Store(name, l)
}

func GetLoader() {

}

type Loader struct {
}

type Importer interface {
}

type importer struct {
	IntPropertyMap    map[int]int
	FloatPropertyMap  map[int]float64
	StringProperties  map[int]string
	MatrixPropertyMap map[int]int
	PointerProperties map[int]*common.AiMatrix4x4
}

func NewImporter() Importer {
	return &importer{}
}

func (im *importer) ReadFile(path string) (*core.AiScene, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

}

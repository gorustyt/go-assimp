package assimp

import (
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver"
	"github.com/gorustyt/go-assimp/driver/protoBin"
	"os"
)

// AiImportFile Reads the given file and returns its content.
func ParseFile(path string) (*core.AiScene, error) {
	im := driver.NewImporter()
	s, err := im.ReadFile(path, 0)
	im = nil
	return s, err
}

func ParseToProtoFile(fromPath, toPath string) error {
	_, err := os.Stat(fromPath)
	if err != nil {
		err = os.MkdirAll(fromPath, 0664)
		if err != nil {
			return err
		}
		logger.WarnF("find path:%v error:%v ,default will  mkdir -r", fromPath, err)
	}
	p, err := ParseFile(fromPath)
	if err != nil {
		return err
	}
	return protoBin.WriteProto(toPath, p)
}

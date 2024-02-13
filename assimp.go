package assimp

import (
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
	"github.com/gorustyt/go-assimp/driver/protoBin"
	"os"
)

// AiImportFile Reads the given file and returns its content.
func ParseFile(path string, steps ...iassimp.AiPostProcessSteps) (*core.AiScene, error) {
	flag := iassimp.AiPostProcessSteps(0)
	for _, v := range steps {
		flag = flag.Add(v)
	}
	im := driver.NewImporter()
	s, err := im.ReadFile(path, flag.Flag())
	im = nil
	return s, err
}

func ParseToProtoFile(fromPath, toPath string, steps ...iassimp.AiPostProcessSteps) error {
	flag := iassimp.AiPostProcessSteps(0)
	for _, v := range steps {
		flag = flag.Add(v)
	}
	_, err := os.Stat(fromPath)
	if err != nil {
		err = os.MkdirAll(fromPath, 0664)
		if err != nil {
			return err
		}
		logger.WarnF("find path:%v error:%v ,default will  mkdir -r", fromPath, err)
	}
	p, err := ParseFile(fromPath, flag)
	if err != nil {
		return err
	}
	return protoBin.WriteProto(toPath, p)
}

package iassimp

import (
	"github.com/gorustyt/go-assimp/common/config"
	"github.com/gorustyt/go-assimp/core"
)

type Loader interface {
	CanRead(checkSig bool) bool
	Read(pScene *core.AiScene) (err error)
	InitConfig(cfg *config.Config)
}

type Importer interface {
	ReadFile(path string, pFlags int) (*core.AiScene, error)
}

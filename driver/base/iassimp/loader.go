package iassimp

import (
	"assimp/common/config"
	"assimp/core"
)

type Loader interface {
	CanRead(checkSig bool) bool
	Read(pScene *core.AiScene) (err error)
	InitConfig(cfg *config.Config)
}

type Importer interface {
	ReadFile(path string, pFlags int) (*core.AiScene, error)
}

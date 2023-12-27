package iassimp

import (
	"assimp/core"
)

type Loader interface {
	CanRead(checkSig bool) bool
	Read(pScene *core.AiScene) (err error)
	Close()
}

type Importer interface {
	ReadFile(path string, pFlags int) (*core.AiScene, error)
}

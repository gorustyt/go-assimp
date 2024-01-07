package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*PretransformVertices)(nil)
)

type PretransformVertices struct {
}

func (p PretransformVertices) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_PreTransformVertices) != 0
}

func (p PretransformVertices) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (p PretransformVertices) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

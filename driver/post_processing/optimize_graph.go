package post_processing

import (
	"assimp/core"
	"assimp/driver/base/iassimp"
	"context"
)

var (
	_ iassimp.PostProcessing = (*OptimizeGraphProcess)(nil)
)

type OptimizeGraphProcess struct {
}

func (o OptimizeGraphProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_OptimizeGraph) != 0
}

func (o OptimizeGraphProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (o OptimizeGraphProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

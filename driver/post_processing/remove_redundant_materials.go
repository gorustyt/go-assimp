package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*RemoveRedundantMatsProcess)(nil)
)

type RemoveRedundantMatsProcess struct {
}

func (r RemoveRedundantMatsProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_RemoveRedundantMaterials) != 0
}

func (r RemoveRedundantMatsProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (r RemoveRedundantMatsProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

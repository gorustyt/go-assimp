package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*RemoveVCProcess)(nil)
)

type RemoveVCProcess struct {
}

func (r RemoveVCProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_RemoveComponent) != 0
}

func (r RemoveVCProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (r RemoveVCProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

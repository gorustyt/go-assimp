package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*ScaleProcess)(nil)
)

type ScaleProcess struct {
}

func (s ScaleProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_GlobalScale) != 0
}

func (s ScaleProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (s ScaleProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

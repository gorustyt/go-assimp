package post_processing

import (
	"assimp/core"
	"assimp/driver/base/iassimp"
	"context"
)

var (
	_ iassimp.PostProcessing = (*SplitLargeMeshesProcessTriangle)(nil)
)

type SplitLargeMeshesProcessTriangle struct {
}

func (s SplitLargeMeshesProcessTriangle) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_SplitLargeMeshes) != 0
}

func (s SplitLargeMeshesProcessTriangle) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (s SplitLargeMeshesProcessTriangle) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

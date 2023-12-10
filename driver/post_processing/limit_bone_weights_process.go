package post_processing

import (
	"assimp/core"
	"assimp/driver/base/iassimp"
	"context"
)

var (
	_ iassimp.PostProcessing = (*LimitBoneWeightsProcess)(nil)
)

type LimitBoneWeightsProcess struct {
}

func (l LimitBoneWeightsProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_LimitBoneWeights) != 0
}

func (l LimitBoneWeightsProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (l LimitBoneWeightsProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

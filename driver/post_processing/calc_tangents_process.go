package post_processing

import (
	"assimp/core"
	"assimp/driver/base/iassimp"
	"context"
)

var (
	_ iassimp.PostProcessing = (*CalcTangentsProcess)(nil)
)

type CalcTangentsProcess struct {
	configMaxAngle float64
	configSourceUV int
}

func (c CalcTangentsProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (c CalcTangentsProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_CalcTangentSpace) != 0
}

func (c CalcTangentsProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

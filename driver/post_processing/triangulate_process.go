package post_processing

import (
	"assimp/core"
	"assimp/driver/base/iassimp"
	"context"
)

var (
	_ iassimp.PostProcessing = (*TriangulateProcess)(nil)
)

type TriangulateProcess struct {
}

func (t TriangulateProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_Triangulate) != 0
}

func (t TriangulateProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (t TriangulateProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

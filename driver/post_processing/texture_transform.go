package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*TextureTransformStep)(nil)
)

type TextureTransformStep struct {
}

func (t TextureTransformStep) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_TransformUVCoords) != 0
}

func (t TextureTransformStep) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (t TextureTransformStep) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

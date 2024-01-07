package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*OptimizeMeshesProcess)(nil)
)

const (
	NotSet   = 0xffffffff
	DeadBeef = 0xdeadbeef
)

type OptimizeMeshesProcess struct {
	//! @see EnablePrimitiveTypeSorting
	pts bool

	//! @see SetPreferredMeshSizeLimit
	max_verts, max_faces int
}

func (o OptimizeMeshesProcess) IsActive(pFlags int) bool {
	if 0 != (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_OptimizeMeshes) {
		o.pts = 0 != (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_SortByPType)
		if 0 != (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_SplitLargeMeshes) {
			o.max_verts = DeadBeef
		}
		return true
	}
	return false
}

func (o OptimizeMeshesProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (o OptimizeMeshesProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

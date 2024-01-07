package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*ArmaturePopulate)(nil)
)

type ArmaturePopulate struct {
	iassimp.BasePostProcessing
}

func (a ArmaturePopulate) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (a ArmaturePopulate) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_PopulateArmatureData) != 0
}

func (a ArmaturePopulate) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

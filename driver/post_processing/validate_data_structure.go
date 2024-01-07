package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*ValidateDSProcess)(nil)
)

type ValidateDSProcess struct {
}

func (v ValidateDSProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (v ValidateDSProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_ValidateDataStructure) != 0
}

func (v ValidateDSProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

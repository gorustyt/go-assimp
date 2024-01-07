package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
)

var (
	_ iassimp.PostProcessing = (*MakeVerboseFormatProcess)(nil)
)

type MakeVerboseFormatProcess struct {
}

func (m MakeVerboseFormatProcess) IsActive(pFlags int) bool {
	// NOTE: There is no direct flag that corresponds to
	// this postprocess step.
	return false
}

func (m MakeVerboseFormatProcess) Execute(pScene *core.AiScene) {
	//TODO implement me
	panic("implement me")
}

func (m MakeVerboseFormatProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

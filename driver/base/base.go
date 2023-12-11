package base

import (
	"assimp/common/reader"
	"assimp/driver/base/iassimp"
	"context"
)

type BaseImporter struct {
	Reader        reader.LineReader
	loader        iassimp.Loader
	bExtraVerbose bool
	context       context.Context
}

func (base *BaseImporter) Init(loader iassimp.Loader, reader *reader.AiReader) {
	base.Reader = reader.GetLineReader()
	base.loader = loader
}

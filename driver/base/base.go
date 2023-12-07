package base

import (
	"assimp/common"
	"strings"
)

type BaseImporter struct {
	Reader *common.AiReader
}

// Check for magic bytes at the beginning of the file.

func (base *BaseImporter) CheckMagicToken(magic string) bool {
	return base.Reader.GetLineNum() == 1 && strings.HasPrefix(base.Reader.GetLine(), magic)

}

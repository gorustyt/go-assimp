package common

import "github.com/gorustyt/go-assimp/common/logger"

func AiAssert(ok bool, msg ...string) {
	if !ok {
		logger.InfoF("Ai assert not ok msg %v ", msg)
	}
}

package common

import "assimp/common/logger"

func AiAssert(ok bool, msg ...string) {
	if !ok {
		logger.InfoF("Ai assert not ok msg %v ", msg)
	}
}

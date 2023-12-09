package common

import "assimp/common/logger"

func AiAssert(ok bool) {
	if !ok {
		logger.InfoF("Ai assert not ok")
	}
}

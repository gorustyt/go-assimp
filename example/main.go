package main

import (
	"assimp/common/logger"
	"assimp/driver"
	"go.uber.org/zap"
)

func main() {
	im := driver.NewImporter()
	p, err := im.ReadFile("./exmaple_data/BLEND/4Cubes4Mats_248.blend", 0)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	_ = p
}

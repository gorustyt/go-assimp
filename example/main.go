package main

import (
	"assimp/common/logger"
	"assimp/driver"
	"go.uber.org/zap"
)

func main() {
	im := driver.NewImporter()
	p, err := im.ReadFile("./example_data/models/AC/closedLine.ac")
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	_ = p
}

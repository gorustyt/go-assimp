package awesomeProject2

import (
	"assimp/core"
	"assimp/driver"
)

// AiImportFile Reads the given file and returns its content.
func ParseFile(path string) (*core.AiScene, error) {
	im := driver.NewImporter()
	s, err := im.ReadFile(path, 0)
	im = nil
	return s, err
}

package awesomeProject2

// AiImportFile Reads the given file and returns its content.
func AiImportFile(path string, pFlags int) error {
	return AiImportFileExWithProperties(path, pFlags, nil)
}

func AiImportFileExWithProperties(path string, pFlags int, props []AiPropertyStore) error {

}

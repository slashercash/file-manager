package internal

import (
	"path/filepath"
	"testing"
)

func TestReadDate(t *testing.T) {

	filePath, _ := filepath.Abs("../../file-manager-test.jpg")

	dateStr, _ := ReadDate(filePath)

	println(dateStr)
}

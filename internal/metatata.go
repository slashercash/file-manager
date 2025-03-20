package internal

import (
	"errors"
	"os"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

func ReadDate(filePath string) (date string, err error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if errC := file.Close(); errC != nil && err == nil {
			err = errC
		}
	}()

	x, err := exif.Decode(file)
	if err != nil {
		return "", err
	}

	tag, err := x.Get(exif.DateTimeOriginal)
	if err != nil {
		return "", err
	}

	if tag.Format() != tiff.StringVal {
		return "", errors.New("DateTime[Original] not in string format")
	}

	return strings.TrimRight(string(tag.Val), "\x00"), nil
}

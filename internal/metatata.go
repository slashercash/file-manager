package internal

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/djherbis/times"
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
		fileInfo, errStat := os.Stat(filePath)
		if errStat != nil {
			return "", errStat
		}
		d := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		cTime := time.Unix(0, d.LastWriteTime.Nanoseconds())
		return cTime.Format("2006:01:02"), nil
	}

	tag, err := x.Get(exif.DateTimeOriginal)
	if err != nil {
		t, err := times.Stat(filePath)
		if err != nil || !t.HasBirthTime() {
			return "", err
		}
		return t.BirthTime().Format("2006:01:02"), err
	}

	if tag.Format() != tiff.StringVal {
		return "", errors.New("DateTime[Original] not in string format")
	}

	return strings.TrimRight(string(tag.Val), "\x00"), nil
}

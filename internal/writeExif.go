package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/djherbis/times"
)

func WriteExif(rootDir string) error {
	rootEntries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range rootEntries {
		if entry.IsDir() {
			continue
		}

		splitted := strings.Split(entry.Name(), ".")
		if len(splitted) != 2 {
			continue
		}

		filePath := path.Join(rootDir, entry.Name())

		t, err := times.Stat(filePath)
		if err != nil {
			return err
		}

		var updateArg string
		switch strings.ToLower(splitted[1]) {
		case "mp4":
			updateArg = fmt.Sprintf("-CreateDate=\"%s\"", t.ModTime().In(time.UTC).Format("2006:01:02 15:04:05"))
		case "png":
			updateArg = fmt.Sprintf("-CreationTime=\"%s\"", t.ModTime().Format("2006:01:02 15:04:05"))
		default:
			continue
		}

		cmd := exec.Command("exiftool", "-overwrite_original", updateArg, filePath)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}

	}

	return nil
}

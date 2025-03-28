package main

import (
	"bytes"
	"errors"
	"file-manager/internal"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := subMain(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func subMain(args []string) error {

	if len(args) != 1 {
		return errors.New("expected arguments: [rootDir]")
	}
	rootDir, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	// copyDateToPartnerFile(rootDir)
	// return nil

	return internal.WriteExif(rootDir)

	rootEntries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	var messages []string

	for _, rootEntry := range rootEntries {
		if !rootEntry.IsDir() {
			messages = append(messages, fmt.Sprintf("ROOT/%s: only directories expected on this level", rootEntry.Name()))
			continue
		}
		entries, err := os.ReadDir(filepath.Join(rootDir, rootEntry.Name()))
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				messages = append(messages, fmt.Sprintf("ROOT/%s/%s: only directories expected on this level", rootEntry.Name(), entry.Name()))
				continue
			}
			msgs, err := scanDir(filepath.Join(rootDir, rootEntry.Name(), entry.Name()))
			if err != nil {
				return err
			}
			messages = append(messages, msgs...)
		}
	}

	for _, message := range messages {
		println(message)
	}

	print(len(messages))

	return nil
}

func scanDir(dir string) ([]string, error) {
	yearStr, monthStr, err := yearAndMonthFromDirPath(dir)
	if err != nil {
		return nil, err
	}

	dirStr := fmt.Sprintf("ROOT/%s/%s", yearStr, monthStr)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var messages []string
	for _, entry := range entries {
		if entry.IsDir() {
			messages = append(messages, fmt.Sprintf("%s/%s: only files expected on this level", dirStr, entry.Name()))
			continue
		}
		dateStr, _ := internal.ReadDate(filepath.Join(dir, entry.Name()))
		// if err != nil {
		// 	return nil, err
		// }
		splitted := strings.Split(dateStr, ":")
		if len(splitted) < 2 || yearStr != splitted[0] || monthStr != splitted[1] {
			messages = append(messages, fmt.Sprintf("%s/%s: date does not match with path - %s", dirStr, entry.Name(), dateStr))
		}
	}

	if len(messages) == 0 {
		messages = append(messages, fmt.Sprintf("%s: OK", dirStr))
	}

	return messages, nil
}

func yearAndMonthFromDirPath(dir string) (year string, month string, err error) {
	dir, month = filepath.Split(filepath.Clean(dir))
	if dir != "" {
		_, year = filepath.Split(filepath.Clean(dir))
	} else {
		return "", "", errors.New("could not get year and month from path")
	}
	return year, month, nil
}

func copyDateToPartnerFile(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {

		splitted := strings.Split(entry.Name(), ".")

		if entry.IsDir() || len(splitted) != 2 || strings.ToLower(splitted[1]) != "jpg" {
			continue
		}

		pngFile := filepath.Join(dir, fmt.Sprintf("%s.png", splitted[0]))

		buf := bytes.Buffer{}
		cmd := exec.Command("exiftool", "-CreationTime", pngFile)
		cmd.Stdout = &buf
		if err := cmd.Run(); err != nil {
			return
		}
		dateTimeStr := strings.Split(buf.String(), " : ")[1][0:19]

		jpgFile := filepath.Join(dir, entry.Name())
		updateArg := fmt.Sprintf("-DateTimeOriginal=\"%s\"", dateTimeStr)
		cmdUpdate := exec.Command("exiftool", "-overwrite_original", updateArg, jpgFile)
		if err := cmdUpdate.Run(); err != nil {
			return
		}
	}
}

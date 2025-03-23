package main

import (
	"errors"
	"file-manager/internal"
	"fmt"
	"os"
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

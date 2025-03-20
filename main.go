package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	print(len(messages))

	return nil
}

func scanDir(dir string) ([]string, error) {
	yearStr, monthStr, err := yearAndMonthFromDirPath(dir)
	if err != nil {
		return nil, err
	}
	_, err = strconv.Atoi(yearStr)
	if err != nil {
		return nil, err
	}
	_, err = strconv.Atoi(monthStr)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var messages []string
	for _, entry := range entries {
		if entry.IsDir() {
			messages = append(messages, fmt.Sprintf("ROOT/%s/%s/%s: only files expected on this level", yearStr, monthStr, entry.Name()))
		}
	}

	return messages, nil
}

func yearAndMonthFromDirPath(dir string) (year string, month string, err error) {
	dir, month = filepath.Split(filepath.Clean(dir))
	if dir != "" {
		_, year = filepath.Split(filepath.Clean(dir))
	} else {
		return "", "", errors.New("Could not get year and month from path")
	}
	return year, month, nil
}

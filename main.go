package main

import (
	"fmt"
	"os"
)

func main() {
	if err := subMain(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func subMain(args []string) error {
	return nil
}

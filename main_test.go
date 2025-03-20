package main

import "testing"

func TestYearAndMonthFromDirPath(t *testing.T) {

	var dirPath = "c:\\ROOT\\2025\\01"

	year, month, _ := yearAndMonthFromDirPath(dirPath)

	print(year, month)

}

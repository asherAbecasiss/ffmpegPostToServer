package main

import (
	"io/ioutil"
	"os"
)

func CountFileInFolder(path string) int {
	files, _ := ioutil.ReadDir(path)

	return len(files)
}
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

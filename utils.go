package main

import (
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// getFiles returns an array of files names found in a directory.
// To get all files leave the ext blank.
func getFiles(dirPath string, ext string) ([]string, error) {

	var fileNames []string

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fileNames, err
	}

	for i := 0; i < len(files); i++ {

		fn := files[i].Name()

		if files[i].IsDir() {
			continue
		}

		if ext != "" {
			v := strings.Split(fn, ".")
			if ext[1:] != v[len(v)-1] {
				continue
			}
		}

		fileNames = append(fileNames, fn)
	}

	return fileNames, nil
}

// fileOrDirectoryExists checks existance of file or directory.
func fileOrDirectoryExists(path string) bool {
	if path == "" {
		return false
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// randInt gets a random int between two numbers.
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	x := max - min
	// In case the max is less than min, take the absolute value.
	x = int(math.Abs(float64(x)))
	return min + rand.Intn(x)
}

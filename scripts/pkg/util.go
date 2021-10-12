package pkg

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

func CheckError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func GetProjectDirectory() string {
	currentDir := filepath.Dir(GetCurrentDir())
	fmt.Println(filepath.Split(currentDir))
	return filepath.Dir(currentDir)

}

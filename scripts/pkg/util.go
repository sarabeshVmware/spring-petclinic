package pkg

import (
	"log"
	"path/filepath"
)

func CheckError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

func CheckFileExtension(path string, ext string) {
	fileExtension := filepath.Ext(path)
	log.Println("File", path, "ext is", fileExtension)
	if fileExtension != ext {
		log.Fatalln("Please change file extension to: ", ext)
	}
}

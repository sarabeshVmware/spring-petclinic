package pkg

import (
	"log"
)

func CheckError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

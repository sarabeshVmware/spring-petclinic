package docker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func DockerLogin(regirstryServer string, username string, password string) error {
	log.Printf("executing docker login to  %s", regirstryServer)

	// create temporary file for password
	tempFile, err := ioutil.TempFile("", "password.json")
	if err != nil {
		log.Printf("error while creating tempfile for tap values schema")
	} else {
		log.Printf("created tempfile")
	}
	defer os.Remove(tempFile.Name())
	err = os.WriteFile(tempFile.Name(), []byte(password), 0677)
	if err != nil {
		log.Printf("error while writing to file %s", tempFile.Name())
		log.Printf("error: %s", err)
	} else {
		log.Printf("file %s written", tempFile.Name())
	}

	// execute cmd
	cmd := fmt.Sprintf("docker login %s -u %s --password-stdin < %s", regirstryServer, username, tempFile.Name())
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("docker login to %s successfull", regirstryServer)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("error while executing docker login to %s", regirstryServer)
		log.Printf("output: %s", output)
	}

	return err
}

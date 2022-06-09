package pkg

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
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

func ExecuteCmd(command string) (string, error) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	// If argument values have spaces in them and passed within single quotes
	for i, value := range arguments {
		if strings.Contains(value, "'") {
			arguments[i] = value + " " + arguments[i+1]
			arguments = append(arguments[:i+1], arguments[i+2:]...)
		}
	}
	// Bypassing single quotes during execution
	for i, value := range arguments {
		arguments[i] = strings.Replace(value, "'", "", -1)
	}
	cmd := exec.Command(commandName, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("Command executed: %s %s", commandName, strings.Join(arguments, " "))
	if err != nil {
		log.Printf("ERROR : %s", err.Error())
		log.Printf("OUTPUT : %s ", string(stdoutStderr))
	} else {
		log.Printf("Output: %s", string(stdoutStderr))
	}
	// if stdoutStderr contains throttling msg, then remove that line before return
	if strings.Contains(string(stdoutStderr), "due to client-side throttling, not priority and fairness") {
		index := strings.Index(string(stdoutStderr), "\n")
		result := string(stdoutStderr)[index+1 : len(string(stdoutStderr))]
		return result, err
	}
	return string(stdoutStderr), err
}

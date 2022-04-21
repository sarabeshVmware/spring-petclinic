package linux_util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

func stripAnsiEscapeSequence(s string) string {
	// Following matches most of the ANSI escape codes, beyond just colors, including the extended VT100 codes, archaic/proprietary printer codes, etc.
	const ansi = "[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(s, "")
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

func ExecuteCmdNoLog(command string) (string, error) {
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
	if err != nil {
		log.Printf("error: %s", err.Error())
		log.Printf("output: %s", stdoutStderr)
	} else {
		log.Printf("output: %s", string(stdoutStderr))
	}
	return string(stdoutStderr), err
}

func ExecuteCmdInBashMode(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("command executed: %s", command)
	if err != nil {
		log.Printf("error: %s", err.Error())
		log.Printf("output: %s", stdoutStderr)
	} else {
		log.Printf("output: %s", string(stdoutStderr))
	}

	return string(stdoutStderr), err
}

func RunCommandWithOutWait(command string) (*os.Process, error) {
	var proc *os.Process
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	err := cmd.Start()
	proc = cmd.Process
	return proc, err
}

func RunBashFile(filepath string, executefrom string) (string, error) {
	cmd := exec.Command(filepath)
	if executefrom != "" {
		cmd.Dir = executefrom
	}
	stdoutStderr, err := cmd.Output()
	if err != nil {
		log.Printf("error: %s", err.Error())
		log.Printf("output: %s", stdoutStderr)
	} else {
		log.Printf("output: %s", string(stdoutStderr))
	}

	return string(stdoutStderr), err
}

type span struct {
	start int
	end   int
}

func indices(s string) []span {
	s = stripAnsiEscapeSequence(s)
	f := unicode.IsSpace
	spans := make([]span, 0, 32)
	start := -1 // valid span start if >= 0
	for end, rune := range s {
		if f(rune) {
			if start >= 0 {
				spans = append(spans, span{start, end})
				start = ^start
			}
		} else {
			if start < 0 {
				start = end
			}
		}
	}
	// Last field might end at EOF.
	if start >= 0 {
		spans = append(spans, span{start, len(s)})
	}
	return spans
}

func FieldIndices(s string) []span {
	spans := indices(s)
	for index := range spans {
		if index == 0 {
			continue
		}
		spans[index-1].end = spans[index].start
	}
	return spans
}

func FieldIndicesWithSingleSpace(s string) []span {
	spans := indices(s)
	mergedspans := make([]span, 0, 32)
	skip := 0
	for index := range spans {
		if skip > 0 {
			skip -= 1
			continue
		}
		if index == len(spans)-1 {
			mergedspans = append(mergedspans, span{spans[index].start, spans[index].end})
			continue
		}
		if spans[index].end+1 == spans[index+1].start {
			if index+2 <= len(spans)-1 && spans[index+1].end+1 == spans[index+2].start {
				skip += 2
				mergedspans = append(mergedspans, span{spans[index].start, spans[index+2].end})
			} else {
				skip += 1
				mergedspans = append(mergedspans, span{spans[index].start, spans[index+1].end})
			}

		} else {
			mergedspans = append(mergedspans, span{spans[index].start, spans[index].end})
		}
	}
	for index := range mergedspans {
		if index == 0 {
			continue
		}
		mergedspans[index-1].end = mergedspans[index].start
	}

	return mergedspans
}

func GetFields(s string, spans []span) []string {
	s = stripAnsiEscapeSequence(s)
	// Create strings from field indices.
	if len(s) < spans[len(spans)-1].end { // if last few column values are empty - padding string with spaces to the right
		b := fmt.Sprintf("%s%d%s", "%-", spans[len(spans)-1].end, "v")
		s = fmt.Sprintf(b, s)
	}
	if len(s) > spans[len(spans)-1].end { // if column values exceed column header length
		spans[len(spans)-1].end = len(s)
	}
	a := make([]string, len(spans))
	for i, span := range spans {
		a[i] = strings.TrimSpace(s[span.start:span.end])
	}
	return a
}

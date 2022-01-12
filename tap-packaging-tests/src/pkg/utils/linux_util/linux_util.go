package linux_util

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

func stripAnsiEscapeSequence(s string) string {
	// Following matches most of the ANSI escape codes, beyond just colors, including the extended VT100 codes, archaic/proprietary printer codes, etc.
	// https://stackoverflow.com/questions/25245716/remove-all-ansi-colors-styles-from-strings/29497680
	//const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	const ansi = "[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(s, "")
}

func ExecuteCmd(command string) (string, error) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("Command executed: %s %s", commandName, strings.Join(arguments, " "))
	if err != nil {
		log.Printf("ERROR : %s", err.Error())
		log.Printf("OUTPUT : %s", stdoutStderr)
	} else {
		log.Printf("Output: \n%s", string(stdoutStderr))
	}
	return stripAnsiEscapeSequence(string(stdoutStderr)), err
}

type span struct {
	start int
	end   int
}

func FieldIndices(s string) []span {
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
	for index := range spans {
		if index == 0 {
			continue
		}
		spans[index-1].end = spans[index].start
	}
	return spans
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

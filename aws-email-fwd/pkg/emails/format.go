package emails

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"strings"
)

const maxLineLength = 78

// wrapHeaderLine takes a string 'line' as input and wraps it into multiple lines
// based on a maximum line length ('maxLineLength'). It formats the text to fit
// within a certain width, inserting line breaks where necessary. If a semicolon
// is encountered, it attempts to break the line at that point.
// Returns the wrapped lines as a single string with ach line indented by four spaces.
//
// https://datatracker.ietf.org/doc/html/rfc6532#section-3.4
func wrapHeaderLine(line string) string {
	var wrappedLines []string
	words := strings.Fields(line)
	currentLine := ""
	for _, word := range words {
		if len(currentLine) != 0 && len(currentLine)+len(word)+1 > maxLineLength {
			lastSemicolon := strings.LastIndex(currentLine, ";")
			if lastSemicolon >= 0 {
				wrappedLines = append(wrappedLines, currentLine[:lastSemicolon+1])
				currentLine = currentLine[lastSemicolon+1:]
			} else {
				wrappedLines = append(wrappedLines, currentLine)
				currentLine = ""
			}
		}
		if currentLine == "" {
			currentLine = word
		} else {
			currentLine += " " + word
		}
	}
	if currentLine != "" {
		wrappedLines = append(wrappedLines, currentLine)
	}
	return strings.Join(wrappedLines, "\n    ")
}

func FormatEmail(msg *mail.Message) []byte {
	sb := bytes.Buffer{}
	for key, values := range msg.Header {
		for _, value := range values {
			headerLine := fmt.Sprintf("%s: %s\n", key, wrapHeaderLine(value))
			_, err := sb.WriteString(headerLine)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	sb.WriteString("\n")
	bodyData, err := io.ReadAll(msg.Body)
	if err != nil {
		panic(err.Error())
	}
	sb.Write(bodyData)
	return sb.Bytes()
}

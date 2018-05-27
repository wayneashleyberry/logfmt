// Command logfmt is an opinionated log formatter
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/wayneashleyberry/truecolor/pkg/color"
)

type message struct {
	Severity string `json:"severity"`
	Time     string `json:"time"`
	Message  string `json:"message"`
	Caller   string `json:"caller"`
}

var cDebug = color.White().Background(76, 117, 217)
var cInfo = color.White().Background(127, 167, 244)
var cWarn = color.White().Background(235, 155, 63)
var cErr = color.White().Background(222, 134, 77)
var cFatal = color.White().Background(195, 81, 63)

var iconDebug = cDebug.Sprint("[λ]")
var iconInfo = cInfo.Sprint("[i]")
var iconWarn = cWarn.Sprint("[!]")
var iconErr = cErr.Sprint("[‼]")
var iconFatal = cFatal.Sprint("[✝]")

var white = color.Color(255, 255, 255)
var dim = color.Color(140, 140, 140)
var superDim = color.Color(80, 80, 80)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var t time.Time

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			println(t, input)
			break
		}

		time, err := println(t, input)
		if err == nil {
			t = time
		}
	}
}

func println(prev time.Time, input string) (time.Time, error) {
	var msg message
	err := json.Unmarshal([]byte(input), &msg)
	if err != nil {
		fmt.Print(input)
		return time.Now(), err
	}

	t, err := time.Parse("2006-01-02T15:04:05.000-0700", msg.Time)
	if err != nil {
		fmt.Print(input)
		return time.Now(), err
	}

	delta := t.Sub(prev)

	var b strings.Builder

	var fields map[string]interface{}
	err = json.Unmarshal([]byte(input), &fields)

	switch msg.Severity {
	case "debug":
		b.WriteString(iconDebug)
	case "info":
		b.WriteString(iconInfo)
	case "warn":
		b.WriteString(iconWarn)
	case "error":
		b.WriteString(iconErr)
	case "fatal":
		b.WriteString(iconFatal)
	default:
		b.WriteString("[" + msg.Severity + "]")
	}

	b.WriteString(" ")
	b.WriteString(dim.Sprint(t.Format("2006-01-02 15:04:05 MST")))
	b.WriteString(" ")
	b.WriteString(white.Sprint(msg.Message))
	b.WriteString(superDim.Sprint(" [" + msg.Caller + "]"))

	var stacktrace string

	for k, v := range fields {
		m, ok := v.(map[string]interface{})
		if ok {
			for kk, vv := range m {
				combined := k + "." + kk
				fields[combined] = vv
			}
		}
	}

	if !prev.IsZero() {
		b.WriteString(superDim.Sprint(" Δ=" + delta.String()))
	}

	for k, v := range fields {
		if k == "severity" || k == "time" || k == "message" || k == "caller" {
			continue
		}
		if k == "stacktrace" {
			stacktrace = v.(string)
			continue
		}

		strval, stringok := v.(string)
		if stringok {
			b.WriteString(superDim.Sprint(" " + k + "=" + strval))
		}

		floatval, floatok := v.(float64)
		if floatok {
			b.WriteString(superDim.Sprint(" " + k + "=" + fmt.Sprintf("%.2f", floatval)))
		}

		boolval, boolok := v.(bool)
		if boolok {
			var stringbool = "true"
			if !boolval {
				stringbool = "false"
			}
			b.WriteString(superDim.Sprint(" " + k + "=" + stringbool))
		}
	}

	if stacktrace != "" {
		lines := strings.Split(stacktrace, "\n")
		for _, line := range lines {
			b.WriteString(superDim.Sprint("\n|>  " + line))
		}
	}

	fmt.Println(b.String())

	return t, nil
}

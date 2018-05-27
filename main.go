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

func main() {
	reader := bufio.NewReader(os.Stdin)

	cDebug := color.White().Background(76, 117, 217)
	cInfo := color.White().Background(127, 167, 244)
	cWarn := color.White().Background(235, 155, 63)
	cErr := color.White().Background(222, 134, 77)
	cFatal := color.White().Background(195, 81, 63)

	iconDebug := cDebug.Sprint("[λ]")
	iconInfo := cInfo.Sprint("[i]")
	iconWarn := cWarn.Sprint("[!]")
	iconErr := cErr.Sprint("[‼]")
	iconFatal := cFatal.Sprint("[✝]")

	white := color.Color(255, 255, 255)
	dim := color.Color(140, 140, 140)
	superDim := color.Color(80, 80, 80)

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		var msg message
		err = json.Unmarshal([]byte(input), &msg)
		if err != nil {
			fmt.Print(input)
			continue
		}

		t, err := time.Parse("2006-01-02T15:04:05.000-0700", msg.Time)
		if err != nil {
			fmt.Print(input)
			continue
		}

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
			b.WriteString(superDim.Sprint("\n" + stacktrace))
		}

		fmt.Println(b.String())
	}
}

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
			fmt.Println(input)
			continue
		}

		t, err := time.Parse("2006-01-02T15:04:05.000-0700", msg.Time)
		if err != nil {
			fmt.Println(input)
			continue
		}

		var b strings.Builder

		var fields map[string]string
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
		b.WriteString(dim.Sprint(t.String()))
		b.WriteString(" ")
		b.WriteString(white.Sprint(msg.Message))

		for k, v := range fields {
			if k == "severity" || k == "time" || k == "message" {
				continue
			}
			b.WriteString(superDim.Sprint(" " + k + "=" + v))
		}

		fmt.Println(b.String())
	}
}

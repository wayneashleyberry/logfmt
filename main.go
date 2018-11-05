// Command logfmt is an opinionated log formatter
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/wayneashleyberry/truecolor/pkg/color"
)

// HTTPPayload does things
type HTTPPayload struct {
	RequestMethod                  string `json:"requestMethod"`
	RequestURL                     string `json:"requestUrl"`
	RequestSize                    string `json:"requestSize"`
	Status                         int    `json:"status"`
	ResponseSize                   string `json:"responseSize"`
	UserAgent                      string `json:"userAgent"`
	RemoteIP                       string `json:"remoteIp"`
	ServerIP                       string `json:"serverIp"`
	Referer                        string `json:"referer"`
	Latency                        string `json:"latency"`
	CacheLookup                    bool   `json:"cacheLookup"`
	CacheHit                       bool   `json:"cacheHit"`
	CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer"`
	CacheFillBytes                 string `json:"cacheFillBytes"`
	Protocol                       string `json:"protocol"`
}

type message struct {
	Severity    string      `json:"severity"`
	Time        string      `json:"time"`
	Message     string      `json:"message"`
	Caller      string      `json:"caller"`
	HTTPPayload HTTPPayload `json:"httpRequest"`
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
			_, _ = println(t, input)
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
		// UTC (Zulu) time
		t2, err := time.Parse("2006-01-02T15:04:05.000Z", msg.Time)
		if err == nil {
			t = t2
		} else {
			fmt.Print(input)
			return time.Now(), err
		}
	}

	delta := t.Sub(prev)

	var b strings.Builder

	var fields map[string]interface{}
	err = json.Unmarshal([]byte(input), &fields)
	if err != nil {
		fmt.Print(input)
		return time.Now(), err
	}

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
	if msg.HTTPPayload.RequestURL != "" {
		box := color.Black().Background(237, 237, 237)
		b.WriteString(box.Sprint(msg.HTTPPayload.RequestMethod))
		b.WriteString(" ")
		b.WriteString(box.Sprint(fmt.Sprintf("%d", msg.HTTPPayload.Status)))

		if msg.HTTPPayload.ResponseSize != "" {
			b.WriteString(" ")
			size, _ := strconv.ParseInt(msg.HTTPPayload.ResponseSize, 10, 64)
			b.WriteString(box.Sprint(humanize.Bytes(uint64(size))))
		}

		if msg.HTTPPayload.Latency != "" {
			d, err := time.ParseDuration(msg.HTTPPayload.Latency)
			if err == nil {
				b.WriteString(" ")
				b.WriteString(box.Sprint(fmt.Sprintf("%s", d)))
			}
		}

		if msg.HTTPPayload.UserAgent != "" {
			b.WriteString(" ")
			if len(msg.HTTPPayload.UserAgent) > 10 {
				b.WriteString(box.Sprint(msg.HTTPPayload.UserAgent[0:10]))
			} else {
				b.WriteString(box.Sprint(msg.HTTPPayload.UserAgent))
			}
		}

		b.WriteString(" ")
		b.WriteString(white.Sprint(msg.HTTPPayload.RequestURL))
	} else {
		b.WriteString(white.Sprint(msg.Message))
	}
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

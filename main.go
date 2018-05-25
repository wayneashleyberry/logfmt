package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type message struct {
	Severity string `json:"severity"`
	Time     string `json:"time"`
	Caller   string `json:"caller"`
	Message  string `json:"message"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
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

		var fields map[string]string
		err = json.Unmarshal([]byte(input), &fields)

		fmt.Printf("%s [%s] - %s (%s)\n", t, msg.Severity, msg.Message, msg.Caller)
	}
}

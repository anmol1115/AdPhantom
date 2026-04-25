package config

import (
	"bytes"
	"fmt"
	"os"
)

func LoadDefaultUpstream() (string, error) {
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return "", err
	}

	var upstream string
	for line := range bytes.Lines(data) {
		if bytes.Contains(line, []byte("nameserver")) {
			row := bytes.Split(line, []byte(" "))
			if len(row) < 2 {
				return "", fmt.Errorf("Error parsing nameserver for default upstream")
			}

			upstream = string(row[1])
			break
		}
	}

	return upstream, err
}

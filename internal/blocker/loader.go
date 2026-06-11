package blocker

import (
	"bytes"
	"errors"
	"os"
)

func LoadFilterRules(path string) error {
	pathInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		return errors.New("Filter list path should be an accessible directory")
	}

	filterFiles, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range filterFiles {
		err := parseFile(file.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func parseFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for line := range bytes.SplitSeq(data, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte{'#'}) {
			continue
		}
	}

	return nil
}

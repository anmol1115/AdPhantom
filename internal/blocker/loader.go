package blocker

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
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
		_, err := parseFile(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func parseFile(path string) ([]*Rule, error) {
	var rules []*Rule

	data, err := os.ReadFile(path)
	if err != nil {
		return rules, err
	}

	for line := range bytes.SplitSeq(data, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if bytes.Contains(line, []byte{'$'}) {
			index := bytes.Index(line, []byte{'$'})
			line = line[:index]
		}
		if len(line) == 0 || bytes.HasPrefix(line, []byte{'#'}) || bytes.HasPrefix(line, []byte{'!'}) || bytes.Contains(line, []byte{'/'}) {
			continue
		}

		if bytes.HasPrefix(line, []byte{'|', '|'}) {
			rule := parseBlock(line)
			if rule != nil {
				rules = append(rules, rule)
			}
		} else if bytes.HasPrefix(line, []byte{'@', '@'}) {
			rule := parseAllow(line)
			if rule != nil {
				rules = append(rules, rule)
			}
		} else if bytes.Contains(line, []byte{' '}) {
			rule := parseExactHost(line)
			if rule != nil {
				rules = append(rules, rule)
			}
		}
	}

	return rules, nil
}

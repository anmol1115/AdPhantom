package blocker

import (
	"bytes"
	"net"
)

type RuleType int

const (
	ExactBlock RuleType = iota
	ExactHost
	WildcardBlock
	ExactAllow
	WildcardAllow
)

type Rule struct {
	Type     RuleType
	Domain   string
	Pattern  string
	Shortcut string
}

func parseBlock(line []byte) *Rule {
	domain := line[2:]
	domain = bytes.TrimSuffix(domain, []byte{'^'})

	var shortcut []byte
	var rule *Rule

	if bytes.Contains(domain, []byte{'*'}) {
		shortcut = getShortcut(domain)
	}

	if len(shortcut) > 0 {
		rule = &Rule{
			Type:     WildcardBlock,
			Shortcut: string(shortcut),
			Pattern:  string(domain),
		}
	} else {
		rule = &Rule{
			Type:   ExactBlock,
			Domain: string(domain),
		}
	}
	return rule
}

func parseAllow(line []byte) *Rule {
	domain := line[2:]
	domain = bytes.TrimSuffix(domain, []byte{'^'})

	if bytes.HasPrefix(domain, []byte{'|', '|'}) {
		domain = domain[2:]
	}
	var shortcut []byte
	var rule *Rule

	if bytes.Contains(domain, []byte{'*'}) {
		shortcut = getShortcut(domain)
	}

	if len(shortcut) > 0 {
		rule = &Rule{
			Type:     WildcardAllow,
			Shortcut: string(shortcut),
			Pattern:  string(domain),
		}
	} else {
		rule = &Rule{
			Type:   ExactAllow,
			Domain: string(domain),
		}
	}
	return rule
}

func parseExactHost(line []byte) *Rule {
	// Exact host rule does not implement ip usage, it will block all exact hosts regardless
	parts := bytes.Split(line, []byte{' '})

	if len(parts) == 2 && len(parts[0]) > 0 && net.ParseIP(string(parts[0])).To4() != nil {
		rule := &Rule{
			Type:   ExactHost,
			Domain: string(parts[1]),
		}
		return rule
	}
	return nil
}

func getShortcut(domain []byte) []byte {
	parts := bytes.Split(domain, []byte{'*'})
	longestpart := parts[0]

	for _, part := range parts {
		if len(part) > len(longestpart) {
			longestpart = part
		}
	}

	if len(longestpart) < 5 {
		return nil
	}

	return longestpart[:5]
}

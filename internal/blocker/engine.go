package blocker

import (
	"bytes"
	"net"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
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

type FilterList struct {
	exactRule    map[string]*Rule
	wildcardRule map[string][]*Rule
	fallbackRule []*Rule
}

func (fl *FilterList) Match(query string) *Rule {
	query = strings.ToLower(strings.TrimRight(query, "."))

	// Exact match
	if rule := fl.exactQueryMatch(query); rule != nil {
		return rule
	}

	// Wildcard Match
	if rule := fl.wildcardQueryMatch(query); rule != nil {
		return rule
	}

	return fl.fallbackQueryMatch(query)
}

func (fl *FilterList) exactQueryMatch(query string) *Rule {
	tld_plus_one, _ := publicsuffix.EffectiveTLDPlusOne(query)
	sample_query := query
	for {
		if rule, found := fl.exactRule[sample_query]; found {
			return rule
		}
		index := strings.Index(sample_query, ".")
		if sample_query == tld_plus_one || index == -1 {
			break
		}
		sample_query = sample_query[index+1:]
	}

	return nil
}

func (fl *FilterList) wildcardQueryMatch(query string) *Rule {
	for shortcut, rules := range fl.wildcardRule {
		if strings.Contains(query, shortcut) {
			for _, rule := range rules {
				matched, err := path.Match(rule.Pattern, query)
				if err != nil {
					continue
				}
				if matched {
					return rule
				}
			}
		}
	}

	return nil
}

func (fl *FilterList) fallbackQueryMatch(query string) *Rule {
	for _, rule := range fl.fallbackRule {
		matched, err := path.Match(rule.Pattern, query)
		if err != nil {
			continue
		}
		if matched {
			return rule
		}
	}
	return nil
}

func parseBlock(line []byte) *Rule {
	domain := line[2:]
	domain = bytes.TrimSuffix(domain, []byte{'^'})

	var shortcut []byte
	var rule *Rule

	if bytes.Contains(domain, []byte{'*'}) {
		shortcut = getShortcut(domain)
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

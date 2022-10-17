package mach

import (
	"strings"
)

func Run(r string) *Runner {
	parts := quoteAwareSplit(string(r))
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	return NewRunner(cmd{
		cmd:  parts[0],
		args: args,
	})
}

func quoteAwareSplit(str string) []string {
	var parts []string
	var lastQuote rune
	var currString strings.Builder

	for i, r := range str {
		switch r {
		case '"', '\'':
			if i == len(str)-1 {
				currString.WriteRune(r)
				continue
			}

			if lastQuote == 0 {
				lastQuote = r
				continue
			}

			if lastQuote == r {
				parts = append(parts, currString.String())
				currString.Reset()
				lastQuote = 0
			} else {
				currString.WriteRune(r)
			}
		case ' ':
			if lastQuote == '\'' || lastQuote == '"' {
				currString.WriteRune(r)
			} else if currString.Len() != 0 {
				parts = append(parts, currString.String())
				currString.Reset()
			}
		default:
			currString.WriteRune(r)
		}
	}

	if currString.Len() != 0 {
		parts = append(parts, currString.String())
	}

	return parts
}

package templateutil

import (
	"strings"
	"unicode"
)

// Names holds all normalized forms of an input resource name.
type Names struct {
	Pascal      string // "OrderItem"
	Snake       string // "order_item"
	Kebab       string // "order-item"
	PluralKebab string // "order-items"
}

// NormalizeName converts any casing style into all required forms.
func NormalizeName(input string) Names {
	words := splitWords(input)
	pascal := toPascal(words)
	snake := strings.Join(words, "_")
	kebab := strings.Join(words, "-")
	pluralKebab := pluralize(kebab)
	return Names{
		Pascal:      pascal,
		Snake:       snake,
		Kebab:       kebab,
		PluralKebab: pluralKebab,
	}
}

// splitWords breaks an identifier into lowercase word parts.
// Handles camelCase, PascalCase, kebab-case, snake_case.
func splitWords(s string) []string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 && s[i-1] != ' ' {
			b.WriteRune(' ')
		}
		b.WriteRune(unicode.ToLower(r))
	}

	return strings.Fields(b.String())
}

func toPascal(words []string) string {
	var b strings.Builder
	for _, w := range words {
		if len(w) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(w[:1]) + w[1:])
	}
	return b.String()
}

// pluralize appends "s" or "es" — sufficient for simple resource names.
func pluralize(s string) string {
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	return s + "s"
}

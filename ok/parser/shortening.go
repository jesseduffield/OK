package parser

import (
	"strings"
)

func abbreviateUndescoredIdentifier(identifier string) string {
	suggested := strings.ToLower(identifier[0:1])
	for i, char := range identifier {
		if char == '_' && i < len(identifier)-1 {
			suggested += strings.ToLower(string(identifier[i+1]))
		}
	}
	return suggested
}

func abbreviateCamelCasedIdentifier(identifier string) string {
	suggested := strings.ToLower(identifier[0:1])
	for _, char := range identifier[1:] {
		strChar := string(char)
		if strings.ToUpper(strChar) == strChar {
			suggested += strings.ToLower(strChar)
		}
	}

	return suggested
}

// expects an identifier that's downcased with no underscores
func smartShorten(identifier string) string {
	charactersToRemove := len(identifier) - MAX_IDENTIFIER_LENGTH

	isVowel := func(c rune) bool { return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' }

	suggested := identifier[0:1]
	for _, char := range identifier[1:] {
		if isVowel(char) && charactersToRemove > 0 {
			charactersToRemove--
		} else {
			suggested += string(char)
		}
	}

	if charactersToRemove == 0 {
		return suggested
	}

	oldSuggested := suggested
	suggested = suggested[0:1]

	if charactersToRemove > 0 {
		// remove every second letter if still too long
		for i, char := range oldSuggested[1:] {
			if i%2 == 1 && charactersToRemove > 0 {
				charactersToRemove--
			} else {
				suggested += string(char)
			}
		}
	}

	if charactersToRemove > 0 {
		suggested = suggested[0:MAX_IDENTIFIER_LENGTH]
	}

	return suggested
}

func shortenedIdentifier(identifier string) string {
	charactersToRemove := len(identifier) - MAX_IDENTIFIER_LENGTH
	if charactersToRemove <= 0 || strings.Count(identifier, "_") >= charactersToRemove {
		return strings.ToLower(removeUnderscores(identifier))
	}

	wordCount := strings.Count(identifier, "_") + 1
	if wordCount > 2 {
		return abbreviateUndescoredIdentifier(identifier)
	}

	if strings.ToLower(identifier) != identifier && strings.ToUpper(identifier) != identifier {
		// has a mix of lowercase and uppercase letters so must be using camelCase

		wordCount := 1
		for _, char := range identifier[1:] {
			if strings.ToUpper(string(char)) == string(char) {
				wordCount++
			}
		}
		if wordCount > 2 {
			return abbreviateCamelCasedIdentifier(identifier)
		}
	}

	return smartShorten(strings.ToLower(removeUnderscores(identifier)))
}

func removeUnderscores(s string) string { return strings.Replace(s, "_", "", -1) }

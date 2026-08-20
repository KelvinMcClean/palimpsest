package normalize

import "strings"

const specialChars = 			"—–-:.,\"()[]{}\\?"
const specialCharsNoSpace = 			"'"
const specialCharsAuthorTrim = 	"—–-:,\"()[]{}\\?"

func Normalize(title string) string {
	// Implement normalization logic here
	// For example, you might want to trim whitespace, convert to lowercase, etc.
	normalized := strings.TrimSpace(title)
	normalized = strings.ToLower(normalized)
	// Add more normalization rules as needed
	for _, char := range specialChars {
		normalized = strings.ReplaceAll(normalized, string(char), " ")
	}
	for _, char := range specialCharsNoSpace {
		normalized = strings.ReplaceAll(normalized, string(char), "")
	}
	normalized = standardizeSpaces(normalized)
	return normalized
}

func standardizeSpaces(s string) string {
    return strings.Join(strings.Fields(s), " ")
}


func splitAndTrim(s string, chars string) string {
	s = strings.TrimSpace(s)
	idx := strings.IndexAny(s, chars)
	if idx != -1 {
		s = s[:idx]
	}
	return s
}

func NormalizeBase(title string) string {
	title = splitAndTrim(title, specialChars)
	return Normalize(title)
}

func NormalizeAuthorName(name string) string {
	normalized := strings.TrimSpace(name)
	normalized = strings.ToLower(normalized)
	for _, char := range specialChars {
		normalized = strings.ReplaceAll(normalized, string(char), " ")
	}

	normalized = standardizeSpaces(normalized)

	// If initials are present, collapse them into a single string without spaces
	// Eg J. R. R. Tolkien -> jrr tolkien, F. Scott Fitzgerald -> f scott fitzgerald, Tolkein, J. R. R. -> tolkien jrr
	charChecker := strings.Split(normalized, " ")
	// Find all adjacent single-character words and collapse them into a single string without spaces
	var finalChars string
	inSingleCharSequence := true
	for i := 0; i < len(charChecker); i++ {
		if len(charChecker[i]) == 1 {
			if (!inSingleCharSequence && i > 0) {
				finalChars += " "
			}
			finalChars += charChecker[i]
			inSingleCharSequence = true
		} else {
			if (i > 0) {
				finalChars += " "
			}
			inSingleCharSequence = false
			finalChars += charChecker[i]
		}
	}
	
	finalChars = standardizeSpaces(finalChars)
	return finalChars
}

func NormalizeBaseAuthorName(name string) string {
	name = splitAndTrim(name, specialCharsAuthorTrim)
	return NormalizeAuthorName(name)
}
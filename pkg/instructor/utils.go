package instructor

import (
	"strings"
)

func toPtr[T any](val T) *T {
	return &val
}

func prepend[T any](to []T, from T) []T {
	return append([]T{from}, to...)
}

func findMatchingBracket(json *string, start int) int {
	stack := []int{}
	openBracket := rune('{')
	closeBracket := rune('}')

	for i := start; i < len(*json); i++ {
		if rune((*json)[i]) == openBracket {
			stack = append(stack, i)
		} else if rune((*json)[i]) == closeBracket {
			if len(stack) == 0 {
				return -1 // Unbalanced brackets
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return i // Found the matching bracket
			}
		}
	}

	return -1 // Unbalanced brackets
}

func getFirstFullJSONElement(json *string) (element string, remaining string) {
	// Find the index of the matching bracket for the first element.
	matchingBracketIdx := findMatchingBracket(json, 0)
	if matchingBracketIdx == -1 {
		return "", *json // No valid JSON element found
	}

	// Extract the full element (including the matching bracket)
	element = (*json)[:matchingBracketIdx+1]

	// Calculate the remaining string after the element
	remaining = (*json)[matchingBracketIdx+1:]

	// If there's a comma after the element, skip it
	if len(remaining) > 0 && remaining[0] == ',' {
		remaining = remaining[1:] // Remove the comma and continue
	}
	element = strings.TrimLeft(element, "[")

	return element, remaining
}

// Removes any prefixes before the JSON (like "Sure, here you go:")
func trimPrefixBeforeJSON(json *string) string {
	startObject := strings.IndexByte(*json, '{')
	startArray := strings.IndexByte(*json, '[')

	var start int
	if startObject == -1 && startArray == -1 {
		return *json // No opening brace or bracket found, return the original string
	} else if startObject == -1 {
		start = startArray
	} else if startArray == -1 {
		start = startObject
	} else {
		start = min(startObject, startArray)
	}

	return (*json)[start:]
}

// Removes any postfixes after the JSON
func trimPostfixAfterJSON(jsonStr string) string {
	endObject := strings.LastIndexByte(jsonStr, '}')
	endArray := strings.LastIndexByte(jsonStr, ']')

	var end int
	if endObject == -1 && endArray == -1 {
		return jsonStr // No closing brace or bracket found, return the original string
	} else if endObject == -1 {
		end = endArray
	} else if endArray == -1 {
		end = endObject
	} else {
		end = max(endObject, endArray)
	}

	return jsonStr[:end+1]
}

// Extracts the JSON by trimming prefixes and postfixes
func extractJSON(json *string) string {
	trimmedPrefix := trimPrefixBeforeJSON(json)
	trimmedJSON := trimPostfixAfterJSON(trimmedPrefix)
	return trimmedJSON
}

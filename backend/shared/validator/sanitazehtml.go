package validator

import (
	"regexp"
	"strings"
)

// SanitizeHTML removes all tags except basic text formatting
func SanitizeHTML(input string) string {
	// Define allowed tags - only basic text formatting
	allowedTags := map[string]bool{
		"p":          true,
		"b":          true,
		"i":          true,
		"u":          true,
		"strong":     true,
		"em":         true,
		"blockquote": true,
		"ul":         true,
		"ol":         true,
		"li":         true,
		"h1":         true,
		"h2":         true,
		"h3":         true,
		"br":         true,
		"hr":         true,
	}

	// Regular expression to match HTML tags
	re := regexp.MustCompile(`</?([a-zA-Z0-9]+)[^>]*>`)

	// Stack to keep track of open tags
	var tagStack []string
	var sanitized strings.Builder

	// Split input into chunks (tags and text)
	matches := re.FindAllStringSubmatchIndex(input, -1)
	lastIndex := 0

	for _, match := range matches {
		// Append text before the tag
		sanitized.WriteString(input[lastIndex:match[0]])

		// Extract tag name and check if it's a closing tag
		fullTag := input[match[0]:match[1]]
		isClosing := strings.HasPrefix(fullTag, "</")
		tagName := strings.ToLower(input[match[2]:match[3]])

		// Check if tag is allowed
		if !allowedTags[tagName] {
			lastIndex = match[1]
			continue
		}

		if isClosing {
			// Handle closing tag
			if len(tagStack) > 0 && tagStack[len(tagStack)-1] == tagName {
				// Matching close tag found
				tagStack = tagStack[:len(tagStack)-1]
				sanitized.WriteString("</" + tagName + ">")
			}
		} else {
			// Handle opening tag
			if tagName != "br" && tagName != "hr" {
				// Add to stack if it's not self-closing
				tagStack = append(tagStack, tagName)
			}

			// Write the sanitized tag
			if tagName == "br" || tagName == "hr" {
				sanitized.WriteString("<" + tagName + " />")
			} else {
				sanitized.WriteString("<" + tagName + ">")
			}
		}

		lastIndex = match[1]
	}

	// Append remaining text
	sanitized.WriteString(input[lastIndex:])

	// Close any remaining open tags in reverse order
	for i := len(tagStack) - 1; i >= 0; i-- {
		sanitized.WriteString("</" + tagStack[i] + ">")
	}

	return sanitized.String()
}

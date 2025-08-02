package vsf

import (
	"fmt"
	"strings"
)

// computeMaxLengths iterates over all rows and gets the max length
// for each column index.
func computeMaxLengths(rows [][]string) []int {
	var maxLengths []int
	for _, row := range rows {
		for i, cell := range row {
			if i >= len(maxLengths) {
				maxLengths = append(maxLengths, len(cell))
			} else if len(cell) > maxLengths[i] {
				maxLengths[i] = len(cell)
			}
		}
	}
	return maxLengths
}

// ParseLine splits a single line into columns respecting quotes.
// Quoted sections (single or double quotes) are treated as single units
// and delimiters inside quotes are ignored.
//
// Parameters:
//   - line: The input line to parse
//   - delimiter: The delimiter to split on
//
// Returns:
//   - A slice of strings representing the parsed columns
//
// Examples:
//
//	ParseLine("name:age:city", ":")
//	// Returns: ["name", "age", "city"]
//
//	ParseLine("john doe:30:new york", ":")
//	// Returns: ["john doe", "30", "new york"]
//
//	ParseLine("'john:doe':30:'new:york'", ":")
//	// Returns: ["'john:doe'", "30", "'new:york'"]
//	// Delimiters inside quotes are preserved
//
//	ParseLine(`"hello,world",foo,"bar,baz"`, ",")
//	// Returns: [`"hello,world"`, "foo", `"bar,baz"`]
//	// Works with both single and double quotes
func ParseLine(line, delimiter string) []string {
	line = strings.TrimSpace(line)
	var (
		result   []string
		current  strings.Builder
		inQuotes bool
	)

	i := 0
	for i < len(line) {
		char := line[i]

		if char == '"' || char == '\'' {
			inQuotes = !inQuotes
			current.WriteByte(char)
			i++
			continue
		}

		if !inQuotes && strings.HasPrefix(line[i:], delimiter) {
			result = append(result, strings.TrimSpace(current.String()))
			current.Reset()
			i += len(delimiter)
			continue
		}

		current.WriteByte(char)
		i++
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

// Format formats input text by aligning columns based on a delimiter.
// This handles 99% of use cases - just formats everything.
//
// Parameters:
//   - input: The input string to be formatted
//   - delimiter: The delimiter used to split each line into columns
//   - outputDelimiter: The delimiter to use in the output. If empty, uses input delimiter
//
// Returns:
//   - A formatted string with aligned columns
//   - An error if the input is empty
//
// Example:
//
//	Format("name:john\nage:30\ncity:new york", ":", "")
//	// Output: "name : john\nage  : 30\ncity : new york"
func Format(input, delimiter, outputDelimiter string) (string, error) {
	return formatColumns(input, delimiter, outputDelimiter, nil)
}

// FormatWithSeparator formats input text and adds a separator line after the specified line.
// Perfect for adding separators after headers.
//
// Parameters:
//   - input: The input string to be formatted
//   - delimiter: The delimiter used to split each line into columns
//   - outputDelimiter: The delimiter to use in the output. If empty, uses input delimiter
//   - afterLine: Line number (0-based) after which to add separator
//   - sepChar: Character to use for the separator line
//
// Returns:
//   - A formatted string with aligned columns and separator
//   - An error if the input is empty
//
// Example:
//
//	FormatWithSeparator("Index:Directory\n5:/path\n0:/short", ":", "", 0, "-")
//	// Output: "Index : Directory\n------:---------\n5     : /path\n0     : /short"
func FormatWithSeparator(input, delimiter, outputDelimiter string, afterLine int, sepChar string) (string, error) {
	formatted, err := formatColumns(input, delimiter, outputDelimiter, nil)
	if err != nil {
		return "", err
	}

	return insertSeparatorAfterLine(formatted, afterLine, outputDelimiter, sepChar), nil
}

// FormatSkipLines formats input text while skipping certain lines from width calculations.
// Useful when you have existing separator lines or headers that shouldn't affect column widths.
//
// Parameters:
//   - input: The input string to be formatted
//   - delimiter: The delimiter used to split each line into columns
//   - outputDelimiter: The delimiter to use in the output. If empty, uses input delimiter
//   - skipLines: Slice of line numbers (0-based) to skip from width calculations
//
// Returns:
//   - A formatted string with aligned columns
//   - An error if the input is empty
//
// Example:
//
//	FormatSkipLines("name:john\n----:----\nage:30", ":", "", []int{1})
//	// Output: "name : john\n----:----\nage  : 30"
//	// Line 1 (----:----) doesn't affect column widths
func FormatSkipLines(input, delimiter, outputDelimiter string, skipLines []int) (string, error) {
	return formatColumns(input, delimiter, outputDelimiter, skipLines)
}

// formatColumns is the core formatting function used by all public functions
func formatColumns(input, delimiter, outputDelimiter string, skipLines []int) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")

	// Create skip map for O(1) lookup
	skipMap := make(map[int]bool)
	for _, lineNum := range skipLines {
		if lineNum >= 0 && lineNum < len(lines) {
			skipMap[lineNum] = true
		}
	}

	// Parse all lines, but only use non-skipped lines for width calculation
	allSegments := make([][]string, len(lines))
	widthCalcSegments := make([][]string, 0)

	for i, line := range lines {
		parsed := ParseLine(line, delimiter)
		allSegments[i] = parsed

		// Only include in width calculation if not skipped
		if !skipMap[i] {
			widthCalcSegments = append(widthCalcSegments, parsed)
		}
	}

	// Calculate max lengths only from non-skipped lines
	maxLengths := computeMaxLengths(widthCalcSegments)

	if outputDelimiter == "" {
		outputDelimiter = delimiter
	}

	var output strings.Builder
	for i, row := range allSegments {
		if skipMap[i] {
			// For skipped lines, output as-is
			output.WriteString(lines[i])
		} else {
			// For normal lines, format with alignment
			for colIndex, cell := range row {
				output.WriteString(cell)
				if colIndex < len(row)-1 {
					// Ensure we don't go out of bounds on maxLengths
					padding := 0
					if colIndex < len(maxLengths) {
						padding = maxLengths[colIndex] - len(cell)
					}
					output.WriteString(strings.Repeat(" ", padding))
					output.WriteString(" " + outputDelimiter + " ")
				}
			}
		}

		if i < len(allSegments)-1 {
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}

// insertSeparatorAfterLine adds a separator line after the specified line number
func insertSeparatorAfterLine(formatted string, afterLine int, outputDelimiter, sepChar string) string {
	lines := strings.Split(formatted, "\n")

	if afterLine < 0 || afterLine >= len(lines) {
		return formatted // Invalid line number, return unchanged
	}

	// Generate separator based on the target line structure
	separatorLine := generateSeparatorFromLine(lines[afterLine], outputDelimiter, sepChar)

	// Insert separator after specified line
	result := make([]string, 0, len(lines)+1)
	for i, line := range lines {
		result = append(result, line)
		if i == afterLine {
			result = append(result, separatorLine)
		}
	}

	return strings.Join(result, "\n")
}

// generateSeparatorFromLine creates a separator that matches the structure of a formatted line
func generateSeparatorFromLine(formattedLine, outputDelimiter, sepChar string) string {
	if outputDelimiter == "" {
		outputDelimiter = ":"
	}

	delimiterPattern := " " + outputDelimiter + " "
	parts := strings.Split(formattedLine, delimiterPattern)

	var separatorParts []string
	for _, part := range parts {
		separatorParts = append(separatorParts, strings.Repeat(sepChar, len(part)))
	}

	return strings.Join(separatorParts, strings.Repeat(sepChar, 1)+outputDelimiter+strings.Repeat(sepChar, 1))
}

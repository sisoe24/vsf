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

// formatColumns is the internal implementation used by both public functions
func formatColumns(input, delimiter, outputDelimiter string, ignoreHeaderLines int) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")

	if ignoreHeaderLines < 0 {
		ignoreHeaderLines = 0
	}

	if ignoreHeaderLines >= len(lines) {
		// If we're ignoring all lines, just return the input as-is
		return input, nil
	}

	// Split into header and data sections
	headerLines := lines[:ignoreHeaderLines]
	dataLines := lines[ignoreHeaderLines:]

	// Parse only the data lines for alignment calculations
	segments := make([][]string, len(dataLines))
	for i, line := range dataLines {
		parsed := ParseLine(line, delimiter)
		segments[i] = parsed
	}

	// Calculate max lengths only from data rows
	maxLengths := computeMaxLengths(segments)

	if outputDelimiter == "" {
		outputDelimiter = delimiter
	}

	var output strings.Builder

	// Add header lines unchanged
	for _, headerLine := range headerLines {
		output.WriteString(headerLine)
		output.WriteString("\n")
	}

	// Format data lines with alignment
	for _, row := range segments {
		for colIndex, cell := range row {
			output.WriteString(cell)
			if colIndex < len(row)-1 {
				padding := maxLengths[colIndex] - len(cell)
				output.WriteString(strings.Repeat(" ", padding))
				output.WriteString(" " + outputDelimiter + " ")
			}
		}
		output.WriteString("\n")
	}

	return strings.TrimSpace(output.String()), nil
}

// Format formats input text by aligning columns based on a delimiter.
//
// Parameters:
//   - input: The input string to be formatted.
//   - delimiter: The delimiter used to split each line into columns.
//   - outputDelimiter: The delimiter to use in the output. If empty, the input delimiter is used.
//
// Returns:
//   - A formatted string with aligned columns.
//   - An error if the input is empty or if any other error occurs during formatting.
//
// Example:
//
//	Format("name:john\nage:30\ncity:new york", ":", "")
//	// Output:
//	// name : john
//	// age  : 30
//	// city : new york
func Format(input, delimiter, outputDelimiter string) (string, error) {
	return formatColumns(input, delimiter, outputDelimiter, 0)
}

// FormatWithHeader formats input text by aligning columns based on a delimiter.
// Header lines are excluded from width calculations but preserved in output.
//
// Parameters:
//   - input: The input string to be formatted.
//   - delimiter: The delimiter used to split each line into columns.
//   - outputDelimiter: The delimiter to use in the output. If empty, the input delimiter is used.
//   - headerLines: Number of header lines to exclude from alignment calculations.
//
// Returns:
//   - A formatted string with aligned columns.
//   - An error if the input is empty or if any other error occurs during formatting.
//
// Example:
//
//	FormatWithHeader("name,age\njohn,30\namy,25", ",", "|", 1)
//	// Output:
//	// name,age
//	// john | 30
//	// amy  | 25
func FormatWithHeader(input, delimiter, outputDelimiter string, headerLines int) (string, error) {
	return formatColumns(input, delimiter, outputDelimiter, headerLines)
}

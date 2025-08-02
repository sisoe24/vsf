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

// ParseLine splits a line into columns, respecting quotes.
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

// AlignColumns formats input text by aligning columns based on a delimiter.
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
//	AlignColumns("name:john\nage:30\ncity:new york", ":", "")
func AlignColumns(input, delimiter, outputDelimiter string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")

	segments := make([][]string, len(lines))
	for i, line := range lines {
		parsed := ParseLine(line, delimiter)
		segments[i] = parsed
	}

	maxLengths := computeMaxLengths(segments)

	if outputDelimiter == "" {
		outputDelimiter = delimiter
	}

	var output strings.Builder
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

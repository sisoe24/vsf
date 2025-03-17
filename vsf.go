// Package main provides a command-line tool for formatting input text
// by aligning columns based on a specified delimiter.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// computeMaxLengths iterates over all rows and updates the max length
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

// parseLine splits a single line into columns, respecting quotes so that
// delimiters inside quotes are not considered splitting points.
func parseLine(line, delimiter string) []string {
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

// formatInput takes an input string, splits it into lines, and then formats each line
// by aligning the columns based on the specified delimiter.
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
//	formatInput("name:john\nage:30\ncity:new york", ":", "")
func formatInput(input, delimiter, outputDelimiter string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	lines := strings.Split(input, "\n")

	segments := make([][]string, len(lines))
	for i, line := range lines {
		parsed := parseLine(line, delimiter)
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

func main() {
	delimiter := flag.String("d", ":", "Delimiter used. Default ':'")
	outputDelimiter := flag.String("o", "", "Output text with selected delimiter")
	usage := flag.Bool("h", false, "Show usage information")

	flag.Usage = showUsage
	flag.Parse()

	if *usage {
		flag.Usage()
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var input strings.Builder
	for scanner.Scan() {
		input.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	formattedOutput, err := formatInput(input.String(), *delimiter, *outputDelimiter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(formattedOutput)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nDescription:\n")
	fmt.Fprintf(os.Stderr, "  This program formats input text by aligning columns based on a specified delimiter.\n")
	fmt.Fprintf(os.Stderr, "  Input is read from stdin. Each line is split by the delimiter, and columns are padded\n")
	fmt.Fprintf(os.Stderr, "  to align with the widest entry in each column.\n")
	fmt.Fprintf(os.Stderr, "\nExample:\n")
	fmt.Fprintf(os.Stderr, "  echo -e \"name:john\\nage:30\\ncity:new york\" | %s\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  Output:\n")
	fmt.Fprintf(os.Stderr, "    name : john\n")
	fmt.Fprintf(os.Stderr, "    age  : 30\n")
	fmt.Fprintf(os.Stderr, "    city : new york\n")
}

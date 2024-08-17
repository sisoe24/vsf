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
// Example:
//     formatInput("name:john\nage:30\ncity:new york", ":", "")
func formatInput(input, delimiter, outputDelimiter string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return "", fmt.Errorf("empty input")
	}

	var maxLengths []int
	segments := make([][]string, len(lines))

	for i, line := range lines {

		line = strings.TrimSpace(line)
		segments[i] = strings.Split(line, delimiter)

		for j, segment := range segments[i] {
			segment = strings.TrimSpace(segment)
			segments[i][j] = segment

			if j >= len(maxLengths) {
				maxLengths = append(maxLengths, len(segment))
			} else if len(segment) > maxLengths[j] {
				maxLengths[j] = len(segment)
			}
		}

	}

	var output strings.Builder

  if outputDelimiter != "" {
    delimiter = outputDelimiter
  }

	for _, lineSegments := range segments {
		for j, segment := range lineSegments {
			output.WriteString(segment)

			if j < len(lineSegments)-1 {
				padding := strings.Repeat(" ", maxLengths[j]-len(segment))
				output.WriteString(padding + " " + delimiter + " ")
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

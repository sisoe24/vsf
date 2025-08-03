package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sisoe24/vsf"
)

const VERSION = "1.3.0"

func main() {
	var (
		delimiter       = flag.String("d", ":", "Delimiter used.")
		outputDelimiter = flag.String("o", "│", "Output text with selected delimiter")
		sepAfter        = flag.Int("sep-after", -1, "Add separator after this line number (0-based)")
		sepChar         = flag.String("sep-char", "═", "Character to use for separator line")
		skipLines       = flag.String("skip", "", "Comma-separated line numbers to skip from width calculations (0-based)")
		version         = flag.Bool("version", false, "Print current version")
		usage           = flag.Bool("h", false, "Show usage information")
	)

	flag.Usage = showUsage
	flag.Parse()

	if *usage {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Println(VERSION)
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

	inputStr := input.String()
	var formattedOutput string
	var err error

	// Parse skip lines if provided
	var skipLineNumbers []int
	if *skipLines != "" {
		skipLineNumbers, err = parseLineNumbers(*skipLines)
		if err != nil {
			log.Fatalf("Error parsing skip lines: %v", err)
		}
	}

	// Choose the appropriate formatting function
	switch {
	case *sepAfter >= 0:
		// Add separator after specified line
		formattedOutput, err = vsf.FormatWithSeparator(inputStr, *delimiter, *outputDelimiter, *sepAfter, *sepChar)
	case len(skipLineNumbers) > 0:
		// Skip certain lines from width calculations
		formattedOutput, err = vsf.FormatSkipLines(inputStr, *delimiter, *outputDelimiter, skipLineNumbers)
	default:
		// Basic formatting (99% of cases)
		formattedOutput, err = vsf.Format(inputStr, *delimiter, *outputDelimiter)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(formattedOutput)
}

// parseLineNumbers parses comma-separated line numbers like "1,3,5"
func parseLineNumbers(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}

	parts := strings.Split(s, ",")
	lineNumbers := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid line number: %s", part)
		}

		lineNumbers = append(lineNumbers, num)
	}

	return lineNumbers, nil
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "vsf version: %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nDescription:\n")
	fmt.Fprintf(os.Stderr, "  This program formats input text by aligning columns based on a specified delimiter.\n")
	fmt.Fprintf(os.Stderr, "  Input is read from stdin. Choose between basic formatting, adding separators, or\n")
	fmt.Fprintf(os.Stderr, "  skipping certain lines from width calculations.\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  Basic formatting (most common):\n")
	fmt.Fprintf(os.Stderr, "    echo -e \"name:john\\nage:30\\ncity:new york\" | %s\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Output:\n")
	fmt.Fprintf(os.Stderr, "      name : john\n")
	fmt.Fprintf(os.Stderr, "      age  : 30\n")
	fmt.Fprintf(os.Stderr, "      city : new york\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Add separator after header (perfect for fzf):\n")
	fmt.Fprintf(os.Stderr, "    echo -e \"Index:Directory\\n5:/long/path\\n0:/short\" | %s -sep-after 0\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Output:\n")
	fmt.Fprintf(os.Stderr, "      Index : Directory\n")
	fmt.Fprintf(os.Stderr, "      ------:---------\n")
	fmt.Fprintf(os.Stderr, "      5     : /long/path\n")
	fmt.Fprintf(os.Stderr, "      0     : /short\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Skip separator lines from width calculations:\n")
	fmt.Fprintf(os.Stderr, "    echo -e \"name:john\\n----:----\\nage:30\" | %s -skip 1\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Output:\n")
	fmt.Fprintf(os.Stderr, "      name : john\n")
	fmt.Fprintf(os.Stderr, "      ----:----\n")
	fmt.Fprintf(os.Stderr, "      age  : 30\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  CSV with headers:\n")
	fmt.Fprintf(os.Stderr, "    cat data.csv | %s -d ',' -sep-after 0 | fzf --header-lines 2\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Custom separators:\n")
	fmt.Fprintf(os.Stderr, "    echo \"a:b:c\" | %s -o ' | ' -sep-after 0 -sep-char '='\n", os.Args[0])
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sisoe24/vsf"
)

const VERSION = "1.3.0"

func main() {
	var (
		delimiter       = flag.String("d", ":", "Delimiter used. Default ':'")
		outputDelimiter = flag.String("o", "", "Output text with selected delimiter")
		headerLines     = flag.Int("header", 0, "Number of header lines to ignore in alignment calculations")
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

	var formattedOutput string
	var err error

	if *headerLines > 0 {
		formattedOutput, err = vsf.FormatWithHeader(input.String(), *delimiter, *outputDelimiter, *headerLines)
	} else {
		formattedOutput, err = vsf.Format(input.String(), *delimiter, *outputDelimiter)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(formattedOutput)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "vsf version: %s\n", VERSION)
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nDescription:\n")
	fmt.Fprintf(os.Stderr, "  This program formats input text by aligning columns based on a specified delimiter.\n")
	fmt.Fprintf(os.Stderr, "  Input is read from stdin. Each line is split by the delimiter, and columns are padded\n")
	fmt.Fprintf(os.Stderr, "  to align with the widest entry in each column.\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  Basic formatting:\n")
	fmt.Fprintf(os.Stderr, "    echo -e \"name:john\\nage:30\\ncity:new york\" | %s\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    Output:\n")
	fmt.Fprintf(os.Stderr, "      name : john\n")
	fmt.Fprintf(os.Stderr, "      age  : 30\n")
	fmt.Fprintf(os.Stderr, "      city : new york\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  CSV with headers (perfect for fzf):\n")
	fmt.Fprintf(os.Stderr, "    cat data.csv | %s -d ',' -header 1 | fzf --header-lines 1\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Custom output delimiter:\n")
	fmt.Fprintf(os.Stderr, "    echo \"a:b:c\" | %s -o ' | '\n", os.Args[0])
}

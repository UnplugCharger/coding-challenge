package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type config struct {
	filePath       string
	printUsage     bool
	byteCount      bool
	lineCount      bool
	wordCount      bool
	characterCount bool
}

// Construct the usage string for the command.
var usageString = fmt.Sprintf(`Usage: %s <filename> [-b] [-h|--help]
A simple tool resembling 'wc' in functionality.
  -b       Count bytes in the given file.
  -h, --help  Display this help and exit.

Example: 
%s test.txt -b
`, os.Args[0], os.Args[0])

func printUsage(w io.Writer) {
	fmt.Fprintf(w, usageString)
}

// parseArgs interprets command-line arguments.
func parseArgs(args []string) (config, error) {
	c := config{}

	if len(args) == 0 {
		return c, errors.New("include at least one argument")
	}

	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			c.printUsage = true
			return c, nil
		case "-c":
			c.byteCount = true
		case "-l":
			c.lineCount = true
		case "-w":
			c.wordCount = true

		case "-m":
			c.characterCount = true
		default:
			if c.filePath != "" {
				return c, errors.New("multiple filenames provided or unrecognized argument: " + arg)
			}
			c.filePath = arg
		}
	}

	// Ensure the help flag isn't used alongside other flags
	if c.printUsage && (c.byteCount || c.lineCount || c.wordCount) {
		return c, errors.New("cannot mix --help with other flags")
	}

	return c, nil
}

// runCmd handles the logic of counting bytes, words, etc.
func runCmd(r io.Reader, w io.Writer, c config) error {
	if c.printUsage {
		printUsage(w)
		return nil
	}

	// Determine if reading from file or stdin
	var content []byte
	var err error
	if c.filePath != "" {
		content, err = os.ReadFile(c.filePath)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}
	} else {
		content, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %v", err)
		}
	}

	// Execute desired counts
	if c.lineCount {
		fmt.Println(strings.Count(string(content), "\n"), c.filePath)
	}
	if c.byteCount {
		fmt.Println(len(content), c.filePath)
	}
	if c.wordCount {
		fmt.Println(len(strings.Fields(string(content))), c.filePath)
	}
	if c.characterCount {
		str := string(content)
		runeCount := 0

		for range str {
			runeCount++
		}
		fmt.Println(runeCount, c.filePath)
	}

	return nil
}

func main() {
	c, err := parseArgs(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := runCmd(os.Stdin, os.Stdout, c); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

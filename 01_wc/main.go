package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// byteCount

// charCount

// lineCount

//  wordCount

type config struct {
	filePath   string
	printUsage bool
	byteCount  bool
	charCount  bool
	lineCount  bool
	wordCount  bool
}

var usageString = fmt.Sprintf(`Usage: %s <integer> [-h|--help]
A greeter application which prints the name you entered <integer> number
of times.
`, os.Args[0])

func printUsage(w io.Writer) {
	fmt.Fprintf(w, usageString)
}

func parseArgs(args []string) (config, error) {

	var err error

	c := config{}

	if len(args) == 0 {
		return c, errors.New("include at least one argument")
	}

	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			c.printUsage = true
		case "-b":
			c.byteCount = true
		case "-c":
			c.charCount = true
		case "-l":
			c.lineCount = true
		case "-w":
			c.wordCount = true
		default:
			if c.filePath != "" { // If filePath is already set, then there's an issue
				return c, errors.New("multiple filenames provided or unrecognized argument: " + arg)
			}
			c.filePath = arg
			fmt.Println("--------->", c.filePath)
		}

	}
	if c.printUsage && (c.byteCount || c.charCount || c.lineCount || c.wordCount) {
		return c, errors.New("cannot mix --help with other flags")
	}

	return c, err
}
func runCmd(r io.Reader, w io.Writer, c config) error {
	if c.printUsage {
		printUsage(w)
		return nil
	}

	var content []byte
	var err error

	if c.filePath != "" {
		content, err = os.ReadFile(c.filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return errors.New("Error reading file ")
		}
	} else {
		// Reading from standard input
		content, err = io.ReadAll(r)
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return errors.New("error reading from stdin")
		}
	}

	if c.lineCount {
		lines := strings.Count(string(content), "\n")
		fmt.Println(lines)
	}

	if c.byteCount {
		fmt.Println(len(content))
	}

	if c.wordCount {
		words := strings.Fields(string(content))
		fmt.Println(len(words))
	}

	return nil
}

func main() {
	c, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

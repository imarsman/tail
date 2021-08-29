package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
	This app takes a number of lines argument, a "pretty" argument for more
	illustrative output, and a list of paths to files, and for each file gathers
	the number of lines requested, if available, and then prints them out to
	standard out.

	The ideal implementation would use a buffer to read in just enough of each
	file to satisfy the number of lines parameter.
*/

// getLines get lasn num lines in file and return them as a string slice. Return
// an error if for instance a filename is incorrect.
func getLines(num int, head bool, path string) ([]string, int, error) {
	var total int
	file, err := os.Open(path)
	if err != nil {
		return nil, total, err
	}

	// Deferring in case an error occurs
	defer file.Close()

	// A bit inefficient as whole file is read in then out again in reverse
	// order up to num.
	// Since we will have to get the last items we have to read all lines in
	// then shorten the output. Other algorithms would involve avoiding reading
	// all the contents in by using a buffer or counting lines or some other
	// technique.
	var all []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		all = append(all, scanner.Text())
	}
	if scanner.Err() != nil {
		return []string{}, total, scanner.Err()
	}

	total = len(all)

	var lines = make([]string, 0, num)

	if head {
		// Get the first lines instead of the last lines
		if total >= num {
			for i := 0; i < num; i++ {
				lines = append(lines, all[i])
				if len(lines) == num {
					break
				}
			}
		} else {
			for i := 0; i < total; i++ {
				lines = append(lines, all[i])
				if len(lines) == num {
					break
				}
			}
		}
	} else {
		// Get last num lines by iterating backwards
		// Slightly more efficient to pre-allocate capacity to known value.
		for i := len(all) - 1; i > -1; i-- {
			lines = append(lines, all[i])
			if len(lines) == num {
				break
			}
		}

		// Another way to do it, which is easier to follow for me. Sample I found
		// returned the slice but you don't need to do that with a slice when it is
		// not being changed in size. As a rule, though, if the slice might be
		// changed you can pass a pointer to it, though that makes it a bit more
		// cumbersome syntactially. I dealt with it in terms of pointers to
		// experiment with the contorted dereferencing.
		var reverse = func(s *[]string) {
			for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
				(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
			}
		}

		// Call the function just defined
		reverse(&lines)
	}

	return lines, total, nil
}

func printHelp() {
	fmt.Println("Print tail (or head) n lines of one or more files")
	fmt.Println("Example: tail -n 10 file1.txt file2.txt")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var h bool
	// Help flag
	flag.BoolVar(&h, "h", false, "print usage")

	var n int
	// Number of lines to print argument
	flag.IntVar(&n, "n", 10, "number of lines")

	var p bool
	// Pretty printing flag
	flag.BoolVar(&p, "p", false, "add formatting to output")
	flag.BoolVar(&p, "pretty", false, "add formatting to output")

	var printLines bool
	// Pring line numbers flag
	flag.BoolVar(&printLines, "N", false, "show line numbers")

	var head bool
	// Print head lines flag
	flag.BoolVar(&head, "H", false, "print head of file rather than tail")

	flag.Parse()

	if h == true {
		printHelp()
	}

	// If a large amount of processing is required handling output for a file at
	// a time shoud help the garbage collector and memory usage.
	// Added total for more informative output.
	var write = func(fname string, head bool, lines []string, total int) {
		builder := new(strings.Builder)
		headStr := "last"
		if head {
			headStr = "first"
		}
		if p == true {
			builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("-", 50)))
		}
		builder.WriteString(fmt.Sprintf("File %s showing %s %d of %d lines\n", fname, headStr, len(lines), total))
		if p == true {
			builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("-", 50)))
		}
		for i := 0; i < len(lines); i++ {
			if printLines == true {
				builder.WriteString(fmt.Sprintf("%-3d %s\n", i+1, lines[i]))
			} else {
				builder.WriteString(fmt.Sprintf("%s\n", lines[i]))
			}
		}
		fmt.Println(strings.TrimSpace(builder.String()))
	}

	// Iterate through list of files (the bits that are not flags), using a
	// strings builder to prepare output. Strings builder avoids allocation.
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("No files specified. Exiting with usage information")
		fmt.Println()
		printHelp()
	}

	for i := 0; i < len(args); i++ {
		lines, total, err := getLines(n, head, args[i])
		if err != nil {
			// panic if something like a bad filename is used
			panic(err)
		}
		write(args[i], head, lines, total)
	}
}

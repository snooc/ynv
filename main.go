package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
)

var buildVersion = "dev"
var buildCommit = "dev"
var buildDate = "n/a"

var opts struct {
	Files       []string `short:"f" long:"file" description:"A file to process"`
	Interactive bool     `short:"i" long:"interactive" description:"Interactive mode" env:"YNV_INTERACTIVE"`
	Version     bool     `short:"v" long:"version" description:"Print version"`
}

func main() {
	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		default:
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	if opts.Version {
		printVersion()
		os.Exit(0)
	}

	// Store list of files to filter with
	filterFiles := opts.Files

	var files []string
	var buffers []*bytes.Buffer

	// Read from file or stdin
	if len(args) == 2 {
		argFile, err := os.Open(args[1])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer argFile.Close()

		files, buffers, err = parseInput(bufio.NewReader(argFile))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		files, buffers, err = parseInput(bufio.NewReader(os.Stdin))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Interactive mode
	if opts.Interactive {
		in := bytes.Buffer{}
		for _, file := range files {
			in.WriteString(fmt.Sprintf("%s\n", file))
		}

		cmd := exec.Command("fzf", "--ansi", "--no-preview")
		var out bytes.Buffer
		cmd.Stdin = &in
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
		choice := strings.TrimSpace(out.String())
		if choice == "" {
			fmt.Printf("No file selected\n")
			os.Exit(1)
		}

		filterFiles = []string{choice}
	}

	// If no files are specified, print all files
	if len(filterFiles) == 0 {
		for _, file := range files {
			fmt.Printf("%s\n", file)
		}
	} else {
		for i, file := range files {
			for _, f := range filterFiles {
				if f == file {
					buffers[i].WriteTo(os.Stdout)
				}
			}
		}

	}
}

var sourceFileRegex = regexp.MustCompile(`# Source: (.*)`)

func parseInput(in *bufio.Reader) ([]string, []*bytes.Buffer, error) {
	var files []string
	var buffers []*bytes.Buffer

	c := -1

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "---":
			c = c + 1
			files = append(files, fmt.Sprintf("file%d", c))
			buffers = append(buffers, &bytes.Buffer{})
		case sourceFileRegex.MatchString(line):
			matches := sourceFileRegex.FindStringSubmatch(line)
			files[c] = matches[1]
		}

		buffers[c].WriteString(line + "\n")
	}

	return files, buffers, nil
}

func printVersion() {
	fmt.Printf("ynv %s (%s) %s\n", buildVersion, buildCommit, buildDate)
}

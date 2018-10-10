package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

var progname string

type selpgArgs struct {
	startPage, endPage, pageLen, pageType int
	inFilename                            string
	printDest                             string
}

var inputS = flag.IntP("startPage", "s", -1, "(Mandatory) Input Your startPage")
var inputE = flag.IntP("endPage", "e", -1, "(Mandatory) Input Your endPage")
var inputL = flag.IntP("pageLen", "l", 72, "(Optional) Choosing pageLen mode, enter pageLen")
var inputF = flag.BoolP("pageBreak", "f", false, "(Optional) Choosing pageBreaks mode")
var inputD = flag.StringP("printDest", "d", "default", "(Optional) Enter printing destination")

func processArgs(selpg *selpgArgs) {
	if *inputS == -1 || *inputE == -1 {
		fmt.Fprintf(os.Stderr, "%v: --startPage(-s) and --endPage(-e) are necessary\n", progname)
		os.Exit(1)
	}
	// handle mandatory arg
	selpg.startPage = *inputS
	selpg.endPage = *inputE
	selpg.pageLen = *inputL
	if *inputF == true {
		selpg.pageType = 'f'
	}
	selpg.printDest = *inputD
	//there is one more arg
	if flag.NArg() >= 1 {
		if flag.NArg() > 1 {
			fmt.Fprintf(os.Stderr, "%v: You should have one file input\n", progname)
			os.Exit(1)
		}
		_, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		selpg.inFilename = flag.Arg(0)
	}
	switch {
	case selpg.startPage <= 0:
		fmt.Fprintf(os.Stderr, "%v: startPage should be bigger than 0\n", progname)
		os.Exit(1)
	case selpg.endPage < selpg.startPage:
		fmt.Fprintf(os.Stderr, "%v: endPage should be bigger than startPage\n", progname)
		os.Exit(1)
	case selpg.pageLen <= 1:
		fmt.Fprintf(os.Stderr, "%v: pageLen should be bigger than 1\n", progname)
		os.Exit(1)
	case selpg.pageType != 'l' && selpg.pageType != 'f':
		fmt.Fprintf(os.Stderr, "%v: There are only two pageTypes for you to choose: pageLen and pageBreaks\n", progname)
		os.Exit(1)
	}
}

func processInput(selpg *selpgArgs) {
	var inputReader *bufio.Reader
	var outputWriter *bufio.Writer
	var err error
	var cmd *exec.Cmd
	var stdin io.WriteCloser
	var file *os.File
	//set the input source
	if selpg.inFilename == "0" {
		inputReader = bufio.NewReader(os.Stdin)
	} else {
		file, err = os.Open(selpg.inFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		inputReader = bufio.NewReader(file)
	}
	//set the output des
	if selpg.printDest == "default" {
		outputWriter = bufio.NewWriter(os.Stdout)
	} else {
		cmd = exec.Command("lp", "-d", selpg.printDest)
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	//begin two input & output loops
	lineCount, pageCount := 0, 1
	var line []byte
	for {
		if selpg.pageType == 'l' {
			line, err = inputReader.ReadBytes('\n')
		} else {
			line, err = inputReader.ReadBytes('\f')
		}
		if err != nil {
			break
		}
		if selpg.pageType == 'l' {
			lineCount++
			if lineCount > selpg.pageLen {
				lineCount = 1
				pageCount++
			}
		}
		if pageCount >= selpg.startPage && pageCount <= selpg.endPage {
			if selpg.printDest == "default" {
				outputWriter.Write(line)
				outputWriter.Flush()
			} else {
				_, err := io.WriteString(stdin, string(line))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}
		if selpg.pageType == 'f' {
			pageCount++
		}
	}
	if selpg.printDest != "default" {
		stdin.Close()
		stderr, _ := cmd.CombinedOutput()
		fmt.Fprintln(os.Stderr, string(stderr))
	}
}

func main() {
	selpg := selpgArgs{
		startPage:  -1,
		endPage:    -1,
		pageLen:    72,
		pageType:   'l',
		inFilename: "0",
		printDest:  "default",
	}
	progname = os.Args[0]
	flag.Parse()
	processArgs(&selpg)
	processInput(&selpg)
}

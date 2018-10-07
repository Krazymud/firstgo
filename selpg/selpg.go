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

var inputS = flag.Int("s", -1, "(Mandatory) Input Your startPage")
var inputE = flag.Int("e", -1, "(Mandatory) Input Your endPage")
var inputL = flag.Int("l", 72, "(Optional) Choosing pageLen mode, enter pageLen")
var inputF = flag.Bool("f", false, "(Optional) Choosing pageBreaks mode")
var inputD = flag.String("d", "default", "(Optional) Enter printing destination")

func processArgs(selpg *selpgArgs) {
	lenOfa := len(os.Args)
	// check the command-line arguments
	if lenOfa < 5 {
		fmt.Printf("%v: not enough arguments\n", progname)
		flag.Usage()
		os.Exit(1)
	}
	// handle first mandatory arg
	if os.Args[1] != "--s" {
		fmt.Fprintf(os.Stderr, "%v: 1st arg should be --s startPage\n", progname)
		flag.Usage()
		os.Exit(1)
	}
	selpg.startPage = *inputS
	//handle second mandatory arg
	if os.Args[3] != "--e" {
		fmt.Fprintf(os.Stderr, "%v: 2nd arg should be --e endPage\n", progname)
		flag.Usage()
		os.Exit(1)
	}
	selpg.endPage = *inputE
	//now handle optional args
	lsign := false
	fsign := false
	for _, a := range os.Args {
		if a == "--l" {
			lsign = true
			selpg.pageLen = *inputL
		}
		if a == "--f" {
			fsign = true
			selpg.pageType = 'f'
		}
		if a == "--d" {
			selpg.printDest = *inputD
		}
	}
	if lsign && fsign {
		fmt.Fprintf(os.Stderr, "%v: You can only choose one mode: pageLen or pageBreaks\n", progname)
		flag.Usage()
		os.Exit(1)
	}
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
		/*cmd.Stdin = strings.NewReader("sss")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}*/
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	//begin two input & output loops
	if selpg.pageType == 'l' {
		lineCount, pageCount := 0, 1
		var line []byte
		for {
			line, err = inputReader.ReadBytes('\n')
			if err != nil {
				break
			}
			lineCount++
			if lineCount > selpg.pageLen {
				lineCount = 1
				pageCount++
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
		}
	} else {
		pageCount := 1
		var bt []byte
		for {
			bt, err = inputReader.ReadBytes('\f')
			if err != nil {
				if err == io.EOF {
					if selpg.printDest == "default" {
						outputWriter.WriteString(string(bt))
						outputWriter.Flush()
					} else {
						_, err := io.WriteString(stdin, string(bt))
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
					}
				}
				break
			}
			pageCount++
			if pageCount >= selpg.startPage && pageCount <= selpg.endPage {
				if selpg.printDest == "default" {
					outputWriter.WriteString(string(bt))
					outputWriter.Flush()
				} else {
					_, err := io.WriteString(stdin, string(bt))
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}
			}
		}
	}
	if selpg.printDest != "default" {
		stdin.Close()
		stderr, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
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

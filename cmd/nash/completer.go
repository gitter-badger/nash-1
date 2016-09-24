package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/NeowayLabs/nash"
)

func completer(sh *nash.Shell) func(string, int) []string {
	return func(line string, pos int) []string {
		err := sh.Exec("autocomplete", `IFS = ()
nashcompletes <= nash_autocomplete("`+line+`")`)

		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return []string{}
		}

		if value, ok := sh.Getvar("nashcompletes"); ok {
			return []string{value.Str()}
		}

		fmt.Printf("not found\n")

		return []string{}
	}
}

func completer2(sh *nash.Shell) func(string, int) []string {
	return func(line string, pos int) []string {
		var local bool

		line = strings.TrimLeft(line[:pos], " ")

		if len(line) == 0 {
			return []string{}
		}

		completeStr := line

		for i := pos - 1; i >= 0; i-- {
			if line[i] == ' ' {
				completeStr = completeStr[i+1 : pos]
				local = true
				break
			}
		}

		if len(completeStr) > 0 {
			for i := 0; i < len(completeStr); i++ {
				if completeStr[i] == ' ' {
					completeStr = completeStr[:i]
					break
				}
			}
		}

		if (len(completeStr) > 0 && (completeStr[0] == '/' || completeStr[0] == '.')) || local {
			return fzf(line, completeFile(completeStr))
		}

		pathVal := os.Getenv("PATH")
		path := make([]string, 0, 256)

		pathparts := strings.Split(pathVal, ":")
		if len(pathparts) == 1 {
			path = append(path, pathparts[0])
		} else {
			for _, p := range pathparts {
				path = append(path, p)
			}
		}

		return fzf(line, completeInPathList(path, completeStr))
	}
}

func completeInPath(path string, complete string) ([]string, bool) {
	var (
		found bool
	)

	newLine := make([]string, 0, 256)

	if len(complete) == 0 {
		return newLine, found
	}

	files, err := ioutil.ReadDir(path)

	if err != nil {
		return newLine, found
	}

	for _, file := range files {
		fname := file.Name()

		if len(complete) <= len(fname) && strings.HasPrefix(fname, complete) {
			newLine = append(newLine, fname)
			if len(complete) == len(fname) {
				found = true
				break
			}
		}
	}

	return newLine, found
}

func completeInPathList(pathList []string, complete string) []string {
	newLine := make([]string, 0, 256)

	for _, path := range pathList {
		tmpNewLine, found := completeInPath(path, complete)

		if len(tmpNewLine) > 0 {
			newLine = append(newLine, tmpNewLine...)
		}

		if found {
			break
		}
	}

	return newLine
}

func completeCurrentPath(complete string) []string {
	lineStr := complete[2:]
	dirParts := strings.Split(lineStr, "/")
	directory := "./" + strings.Join(dirParts[0:len(dirParts)-1], "/")

	newLine := make([]string, 0, 256)

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		return newLine
	}

	for _, file := range files {
		var cmpStr string

		fname := file.Name()

		if fname == "." {
			continue
		}

		if directory[len(directory)-1] == '/' {
			cmpStr = directory + fname
		} else {
			cmpStr = directory + "/" + fname
		}

		if len(cmpStr) >= len(complete) &&
			strings.HasPrefix(cmpStr, complete) {

			newLine = append(newLine, cmpStr)
		}
	}

	return newLine
}

func completeAbsolutePath(complete string) []string {
	lineStr := complete[1:] // ignore first '/'
	dirParts := strings.Split(lineStr, "/")
	directory := "/" + strings.Join(dirParts[0:len(dirParts)-1], "/")

	newLine := make([]string, 0, 256)

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		return newLine
	}

	for _, file := range files {
		var cmpStr string

		fname := file.Name()

		if directory[len(directory)-1] == '/' {
			cmpStr = directory + fname
		} else {
			cmpStr = directory + "/" + fname
		}

		if len(cmpStr) >= len(complete) && strings.HasPrefix(cmpStr, complete) {
			newLine = append(newLine, cmpStr)
		}
	}

	return newLine
}

func completeFile(complete string) []string {
	llen := len(complete)

	if llen >= 1 {
		if complete[0] != '/' {
			if llen >= 2 && complete[0] == '.' && complete[1] == '/' {
				return completeCurrentPath(complete)
			} else {
				return completeCurrentPath("./" + complete)
			}
		}

		return completeAbsolutePath(complete)

	} else {
		return completeFile("./")
	}
}

func fzf(line string, choices []string) []string {
	cmd := exec.Command("fzf", "-q", line)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return []string{}
	}

	go func() {
		for _, str := range choices {
			fmt.Fprintf(stdin, "%s\n", str)
		}

		stdin.Close()
	}()

	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return []string{}
	}

	return []string{string(out)}
}

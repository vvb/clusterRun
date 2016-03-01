package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var filter = flag.String("filter", "*", "node to filter")
	var fileName = flag.String("file", "commands.cfg", "commands file")
	var prefix = flag.String("cmd", "/usr/bin/vagrant", "vagrant ssh or ssh")
	flag.Parse()

	lines, err := readLines(*fileName)
	if err != nil {
		log.Fatalf("Error reading commands file")
	}

	var node, cmdStr string
	for _, line := range lines {
		if line == "" {
			fmt.Println("skipping empty line..")
			continue
		}

		// remove any preceding spaces
		line = strings.TrimLeft(line, " ")

		// ignore lines that are commented out
		if strings.Index(line, "#") == 0 {
			continue
		}

		// extract node and cmd strings
		array := strings.Split(line, ":")
		node = array[0]
		for i, value := range array {
			if i == 0 {
				continue
			} else if i > 1 {
				cmdStr += ":"
			}
			cmdStr += value
		}

		if !((*filter == "*") || (*filter == node)) {
			continue
		}

		fmt.Println("====>", node, cmdStr)
		cmd := exec.Command(*prefix, "ssh", node, "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		fmt.Println("<====", byteToString(output))
		if err != nil {
			fmt.Println("===ERROR===", err)
		}
		cmdStr = ""
	}
}

func byteToString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

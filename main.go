package main

import (
	"bufio"
	"bytes"
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
		node = strings.Split(line, ":")[0]
		cmdStr = strings.Split(line, ":")[1]

		if !((*filter == "*") || (*filter == node)) {
			continue
		}

		fmt.Println("running on node", node, cmdStr)
		cmd := exec.Command(*prefix, "ssh", node, "-c", cmdStr)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%q\n", out.String())
	}
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

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// openFile open a file for writing
func writeLines(line string, path string) error {

	// overwrite file if it exists
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	_, err = w.WriteString(fmt.Sprintf("%v%v", line, "\n"))
	if err != nil {
		log.Fatal(err)
	}

	// flush outstanding data
	return w.Flush()
}

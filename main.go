package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Print(err)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		data := make([]byte, 8)
		line := ""
		
		for {
				_, err := f.Read(data)
				if err != nil {
					f.Close()
					if errors.Is(err, io.EOF) {
						lines <- line
						os.Exit(0)
					}
					fmt.Print(err)
				}

				
				parts := bytes.Split(data, []byte("\n"))
		
				data = nil
				data = make([]byte, 8)

				for i, v := range parts {
					line += string(v)
					if len(parts) > i + 1 {
						lines <- line
						line = ""
					}
				}
			}
		}()

	return lines
}

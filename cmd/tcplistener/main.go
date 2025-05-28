package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection has been accepted!")

		func(conn net.Conn) {
			lines := getLinesChannel(conn)

			for line := range lines {
				fmt.Printf("%s\n", line)
			}
	
			fmt.Println("Connection Closed.")
			conn.Close()
		} (c)
		
	}	
}

func getLinesChannel(f net.Conn) <-chan string {
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
						close(lines)
						break
					}
					log.Fatal(err)
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

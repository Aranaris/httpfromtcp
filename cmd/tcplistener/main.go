package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/request"
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
			request, err := request.RequestFromReader(conn)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Request line: ")
			fmt.Printf("- Method: %v\n", request.RequestLine.Method)
			fmt.Printf("- Target: %v\n", request.RequestLine.RequestTarget)
			fmt.Printf("- Version: %v\n", request.RequestLine.HttpVersion)

			fmt.Println("Connection Closed.")
			conn.Close()
		} (c)
		
	}	
}

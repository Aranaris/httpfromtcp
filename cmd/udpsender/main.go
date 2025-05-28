package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	u, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
	}

	c, err := net.DialUDP("udp", nil, u)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = c.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
	}
}
